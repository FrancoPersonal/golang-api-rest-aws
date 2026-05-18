package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/dto"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/handlers"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/middleware"
	dynamodbrepo "github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/repositories/dynamodb"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/application/services"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load aws config: %v", err)
	}

	tableName := os.Getenv("CAUCION_TABLE_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	dbClient := awsdynamodb.NewFromConfig(cfg)
	repo := dynamodbrepo.NewSuretyBondRepository(dbClient, tableName)
	service := services.NewSuretyBondService(repo, nil, nil)
	handler := handlers.NewSuretyBondHandler(service)

	logger := log.New(os.Stdout, "", log.LstdFlags)
	list := middleware.Chain(
		handler.ListSuretyBonds,
		middleware.LoggingMiddleware(logger),
		middleware.JWTAuthMiddleware(jwtSecret),
	)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if req.HTTPMethod != "GET" {
			return dto.Fail(404, "not_found", "route not found"), nil
		}
		return list(ctx, req)
	})
}
