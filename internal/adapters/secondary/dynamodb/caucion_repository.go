package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	"github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

// dynamoDBClient is the interface for DynamoDB operations (enables mocking in tests).
type dynamoDBClient interface {
	PutItem(ctx context.Context, params *awsdynamodb.PutItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *awsdynamodb.GetItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error)
	Scan(ctx context.Context, params *awsdynamodb.ScanInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error)
	UpdateItem(ctx context.Context, params *awsdynamodb.UpdateItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, params *awsdynamodb.DeleteItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.DeleteItemOutput, error)
}

// Repository implements secondary.CaucionRepository using Amazon DynamoDB.
type Repository struct {
	client    dynamoDBClient
	tableName string
	log       logger.Logger
}

// New creates a new DynamoDB-backed Repository.
func New(client dynamoDBClient, tableName string, log logger.Logger) *Repository {
	return &Repository{client: client, tableName: tableName, log: log}
}

func (r *Repository) Save(ctx context.Context, c *domain.Caucion) error {
	item, err := attributevalue.MarshalMap(c)
	if err != nil {
		return fmt.Errorf("marshal caucion: %w", err)
	}

	_, err = r.client.PutItem(ctx, &awsdynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put item: %w", err)
	}

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Caucion, error) {
	out, err := r.client.GetItem(ctx, &awsdynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	if out.Item == nil {
		return nil, domain.ErrNotFound
	}

	var c domain.Caucion
	if err := attributevalue.UnmarshalMap(out.Item, &c); err != nil {
		return nil, fmt.Errorf("unmarshal caucion: %w", err)
	}

	return &c, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]domain.Caucion, error) {
	out, err := r.client.Scan(ctx, &awsdynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	cauciones := make([]domain.Caucion, 0, len(out.Items))
	for _, item := range out.Items {
		var c domain.Caucion
		if err := attributevalue.UnmarshalMap(item, &c); err != nil {
			return nil, fmt.Errorf("unmarshal caucion: %w", err)
		}
		cauciones = append(cauciones, c)
	}

	return cauciones, nil
}

func (r *Repository) Update(ctx context.Context, c *domain.Caucion) error {
	item, err := attributevalue.MarshalMap(c)
	if err != nil {
		return fmt.Errorf("marshal caucion: %w", err)
	}

	_, err = r.client.PutItem(ctx, &awsdynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &awsdynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}
