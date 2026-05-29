package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/dto"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/middleware"
	lambdautil "github.com/FrancoPersonal/golang-api-rest-aws/internal/infrastructure/lambda"
)

func main() {
	lambda.Start(newHandler())
}

func newHandler() func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c := lambdautil.InitializeApiHandler()

	create := middleware.Chain(
		c.SuretyBondHandler.CreateSuretyBond,
		middleware.LoggingMiddleware(c.Logger),
		middleware.JWTAuthMiddleware(c.JWTSecret),
	)

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if req.HTTPMethod != "POST" {
			return dto.Fail(404, "not_found", "route not found"), nil
		}
		return create(ctx, req)
	}
}
