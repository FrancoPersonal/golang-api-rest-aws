package dynamodb

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type putItemClientMock struct {
	err    error
	called bool
	items  []map[string]types.AttributeValue
}

func (m *putItemClientMock) PutItem(_ context.Context, _ *awsdynamodb.PutItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error) {
	m.called = true
	return &awsdynamodb.PutItemOutput{}, m.err
}

func (m *putItemClientMock) Scan(_ context.Context, _ *awsdynamodb.ScanInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &awsdynamodb.ScanOutput{Items: m.items}, nil
}

func TestCreateSuccess(t *testing.T) {
	m := &putItemClientMock{}
	repo := NewSuretyBondRepository(m, "suretyBonds")

	err := repo.Create(context.Background(), domain.SuretyBond{ID: "1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	require.NoError(t, err)
	require.True(t, m.called)
}

func TestCreateConflict(t *testing.T) {
	m := &putItemClientMock{err: &types.ConditionalCheckFailedException{}}
	repo := NewSuretyBondRepository(m, "suretyBonds")

	err := repo.Create(context.Background(), domain.SuretyBond{ID: "1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	require.ErrorIs(t, err, domain.ErrConflict)
}

func TestListSuccess(t *testing.T) {
	now := time.Now().UTC()
	item, err := attributevalue.MarshalMap(domain.SuretyBond{ID: "1", Number: "C-1", CreatedAt: now, UpdatedAt: now})
	require.NoError(t, err)

	m := &putItemClientMock{items: []map[string]types.AttributeValue{item}}
	repo := NewSuretyBondRepository(m, "suretyBonds")

	result, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "1", result[0].ID)
}

func TestListInternalError(t *testing.T) {
	m := &putItemClientMock{err: assertiveErr{}}
	repo := NewSuretyBondRepository(m, "suretyBonds")

	result, err := repo.List(context.Background())
	require.Nil(t, result)
	require.ErrorIs(t, err, domain.ErrInternal)
}

func TestCreateMissingTableName(t *testing.T) {
	m := &putItemClientMock{}
	repo := NewSuretyBondRepository(m, "")

	err := repo.Create(context.Background(), domain.SuretyBond{ID: "1"})
	require.ErrorIs(t, err, domain.ErrInternal)
}

func TestCreateDynamoDBInternalError(t *testing.T) {
	m := &putItemClientMock{err: assertiveErr{}}
	repo := NewSuretyBondRepository(m, "suretyBonds")

	err := repo.Create(context.Background(), domain.SuretyBond{ID: "1", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()})
	require.ErrorIs(t, err, domain.ErrInternal)
}

func TestListMissingTableName(t *testing.T) {
	m := &putItemClientMock{}
	repo := NewSuretyBondRepository(m, "")

	result, err := repo.List(context.Background())
	require.Nil(t, result)
	require.ErrorIs(t, err, domain.ErrInternal)
}

func TestListEmptyItems(t *testing.T) {
	m := &putItemClientMock{items: []map[string]types.AttributeValue{}}
	repo := NewSuretyBondRepository(m, "suretyBonds")

	result, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Empty(t, result)
}

type assertiveErr struct{}

func (assertiveErr) Error() string { return "boom" }
