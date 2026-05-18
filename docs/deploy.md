# Deploy

## Stacks

The deployment is split into three independent Serverless stacks:

| Stack | File | Description |
|---|---|---|
| Policies | `serverless.policies.yml` | IAM roles and permissions |
| Infrastructure | `serverless.infrastructure.yml` | DynamoDB table and shared resources |
| API | `serverless.yml` | Lambda functions and API Gateway |

Each stack is independently deployable and the API stack consumes outputs from the other two via CloudFormation.

---

## Deploy commands

```bash
# Review stage variables before deploying
serverless print --stage dev

# Deploy each stack in order
task deploy:policies:dev
task deploy:infrastructure:dev
task deploy:dev

# Or deploy all at once
task deploy:dev    # dev
task deploy:prod   # prod
```

```bash
# Equivalent serverless commands
serverless deploy --config serverless.policies.yml --stage dev
serverless deploy --config serverless.infrastructure.yml --stage dev
serverless deploy --stage dev
```

---

## Environment variables (Lambda)

| Variable | Source | Description |
|---|---|---|
| `CAUCION_TABLE_NAME` | CloudFormation output | DynamoDB table name |
| `JWT_SECRET` | AWS Secrets Manager | JWT signing secret |
| `LOG_LEVEL` | Serverless stage var | Log verbosity: `debug`, `info`, `warn`, `error` (default: `info`) |

---

## Security

- Authentication via OAuth2/JWT — validated in `middleware/auth_jwt.go`
- Secrets stored in AWS Secrets Manager — never hardcoded
- Minimum required IAM permissions defined in `serverless.policies.yml`
- Input validation at the HTTP boundary; business validation in use cases
