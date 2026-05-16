package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type PutItemClient interface {
	PutItem(ctx context.Context, params *awsdynamodb.PutItemInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *awsdynamodb.ScanInput, optFns ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error)
}

type CaucionRepository struct {
	client    PutItemClient
	tableName string
}

func NewCaucionRepository(client PutItemClient, tableName string) *CaucionRepository {
	return &CaucionRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *CaucionRepository) Create(ctx context.Context, caucion domain.Caucion) error {
	if r.tableName == "" || r.client == nil {
		return domain.ErrInternal
	}

	item, err := attributevalue.MarshalMap(caucion)
	if err != nil {
		return domain.ErrInternal
	}

	_, err = r.client.PutItem(ctx, &awsdynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})
	if err != nil {
		var conflict *types.ConditionalCheckFailedException
		if errors.As(err, &conflict) {
			return domain.ErrConflict
		}
		return domain.ErrInternal
	}

	return nil
}

func (r *CaucionRepository) List(ctx context.Context) ([]domain.Caucion, error) {
	if r.tableName == "" || r.client == nil {
		return nil, domain.ErrInternal
	}

	out, err := r.client.Scan(ctx, &awsdynamodb.ScanInput{TableName: aws.String(r.tableName)})
	if err != nil {
		return nil, domain.ErrInternal
	}

	if len(out.Items) == 0 {
		return []domain.Caucion{}, nil
	}

	var cauciones []domain.Caucion
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &cauciones); err != nil {
		return nil, domain.ErrInternal
	}

	return cauciones, nil
}
