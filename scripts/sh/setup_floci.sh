#!/usr/bin/env bash
set -euo pipefail

ENDPOINT=http://localhost:4566
REGION=us-east-1
ACCOUNT_ID=000000000000
TABLE_NAME=cauciones-api-local-cauciones
JWT_SECRET=${JWT_SECRET:-dev-secret}

# All AWS CLI calls run inside the Floci container via awslocal, which is
# pre-wired to http://localhost:4566 with dummy credentials.
CLI="docker exec floci awslocal"

# ---------------------------------------------------------------------------
# Wait for Floci container and its API to be ready
# ---------------------------------------------------------------------------
echo "Waiting for Floci container..."
until docker inspect -f '{{.State.Running}}' floci 2>/dev/null | grep -q true; do
  sleep 1
done

echo "Waiting for Floci API..."
until curl -s "${ENDPOINT}/health" > /dev/null 2>&1 || \
      curl -s "${ENDPOINT}/dynamodb" > /dev/null 2>&1 || \
      $CLI dynamodb list-tables > /dev/null 2>&1; do
  sleep 1
done
echo "Floci is ready"

# ---------------------------------------------------------------------------
# DynamoDB table
# ---------------------------------------------------------------------------
$CLI dynamodb create-table \
  --table-name "$TABLE_NAME" \
  --attribute-definitions AttributeName=id,AttributeType=S \
  --key-schema AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null && echo "DynamoDB table created: $TABLE_NAME" \
  || echo "DynamoDB table already exists: $TABLE_NAME"

# ---------------------------------------------------------------------------
# IAM role (required by Lambda even in local mode)
# ---------------------------------------------------------------------------
ROLE_ARN="arn:aws:iam::${ACCOUNT_ID}:role/lambda-local-role"
$CLI iam create-role \
  --role-name lambda-local-role \
  --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}' \
  2>/dev/null && echo "IAM role created" \
  || echo "IAM role already exists"

# ---------------------------------------------------------------------------
# Lambda functions
# Zip files are mounted at /build/ inside the Floci container.
# ---------------------------------------------------------------------------
ENV_VARS="Variables={CAUCION_TABLE_NAME=${TABLE_NAME},JWT_SECRET=${JWT_SECRET}}"

deploy_lambda() {
  local name=$1
  local zip_container=$2   # path inside the container (under /build/)

  if $CLI lambda get-function --function-name "$name" > /dev/null 2>&1; then
    $CLI lambda update-function-code \
      --function-name "$name" \
      --zip-file "fileb://${zip_container}" > /dev/null
    $CLI lambda update-function-configuration \
      --function-name "$name" \
      --environment "$ENV_VARS" > /dev/null
    echo "Lambda updated: $name"
  else
    $CLI lambda create-function \
      --function-name "$name" \
      --runtime provided.al2023 \
      --role "$ROLE_ARN" \
      --handler bootstrap \
      --zip-file "fileb://${zip_container}" \
      --environment "$ENV_VARS" > /dev/null
    echo "Lambda created: $name"
  fi
}

deploy_lambda post-cauciones-local /build/post_cauciones.zip
deploy_lambda get-cauciones-local  /build/get_cauciones.zip

# ---------------------------------------------------------------------------
# API Gateway REST (v1) — matches events.APIGatewayProxyRequest in handlers
# ---------------------------------------------------------------------------
API_ID=$($CLI apigateway get-rest-apis \
  --query 'items[?name==`cauciones-local`].id' \
  --output text 2>/dev/null)

if [ -z "$API_ID" ] || [ "$API_ID" = "None" ]; then
  API_ID=$($CLI apigateway create-rest-api \
    --name cauciones-local \
    --query 'id' --output text)
  echo "API Gateway created: $API_ID"

  ROOT_ID=$($CLI apigateway get-resources \
    --rest-api-id "$API_ID" \
    --query 'items[?path==`/`].id' --output text)

  RESOURCE_ID=$($CLI apigateway create-resource \
    --rest-api-id "$API_ID" \
    --parent-id "$ROOT_ID" \
    --path-part cauciones \
    --query 'id' --output text)

  POST_ARN="arn:aws:lambda:${REGION}:${ACCOUNT_ID}:function:post-cauciones-local"
  GET_ARN="arn:aws:lambda:${REGION}:${ACCOUNT_ID}:function:get-cauciones-local"
  URI_BASE="arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions"

  $CLI apigateway put-method \
    --rest-api-id "$API_ID" --resource-id "$RESOURCE_ID" \
    --http-method POST --authorization-type NONE > /dev/null
  $CLI apigateway put-integration \
    --rest-api-id "$API_ID" --resource-id "$RESOURCE_ID" \
    --http-method POST --type AWS_PROXY --integration-http-method POST \
    --uri "${URI_BASE}/${POST_ARN}/invocations" > /dev/null

  $CLI apigateway put-method \
    --rest-api-id "$API_ID" --resource-id "$RESOURCE_ID" \
    --http-method GET --authorization-type NONE > /dev/null
  $CLI apigateway put-integration \
    --rest-api-id "$API_ID" --resource-id "$RESOURCE_ID" \
    --http-method GET --type AWS_PROXY --integration-http-method POST \
    --uri "${URI_BASE}/${GET_ARN}/invocations" > /dev/null

  $CLI apigateway create-deployment \
    --rest-api-id "$API_ID" \
    --stage-name local > /dev/null
  echo "API Gateway stage deployed"
else
  echo "API Gateway already exists: $API_ID"
fi

echo ""
echo "Local API ready:"
echo "  POST ${ENDPOINT}/restapis/${API_ID}/local/_user_request_/cauciones"
echo "  GET  ${ENDPOINT}/restapis/${API_ID}/local/_user_request_/cauciones"
