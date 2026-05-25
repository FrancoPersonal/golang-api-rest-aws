package lambda

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/handlers"
	dynamodbrepo "github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/repositories/dynamodb"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/application/services"
)

// APIComponents holds the shared wired dependencies for all Lambda handlers.
type APIComponents struct {
	SuretyBondHandler *handlers.SuretyBondHandler
	JWTSecret         string
	Logger            *log.Logger
}

// InitializeApiHandler loads the AWS config, wires the repository, service and handler,
// and returns the shared components. Must be called once at cold-start.
func InitializeApiHandler() APIComponents {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load aws config: %v", err)
	}

	tableName := os.Getenv("CAUCION_TABLE_NAME")

	dbClient := awsdynamodb.NewFromConfig(cfg)
	repo := dynamodbrepo.NewSuretyBondRepository(dbClient, tableName)
	svc := services.NewSuretyBondService(repo, nil, nil)
	handler := handlers.NewSuretyBondHandler(svc)

	return APIComponents{
		SuretyBondHandler: handler,
		JWTSecret:         os.Getenv("JWT_SECRET"),
		Logger:            log.New(os.Stdout, "", log.LstdFlags),
	}
}
