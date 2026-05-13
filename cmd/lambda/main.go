package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	httphandler "github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/primary/http"
	dynamodbadapter "github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/secondary/dynamodb"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/application"
	"github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log := logger.New(logLevel)

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error("failed to load AWS config", "error", err)
		os.Exit(1)
	}

	tableName := os.Getenv("CAUCION_TABLE_NAME")
	if tableName == "" {
		log.Error("CAUCION_TABLE_NAME environment variable is not set")
		os.Exit(1)
	}

	ddbClient := awsdynamodb.NewFromConfig(cfg)
	repo := dynamodbadapter.New(ddbClient, tableName, log)
	uc := application.New(repo, log)
	router := httphandler.NewRouter(uc, log)

	adapter := httpadapter.NewV2(router)
	lambda.Start(adapter.ProxyWithContext)
}
