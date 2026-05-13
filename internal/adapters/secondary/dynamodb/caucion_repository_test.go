package dynamodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dynamodbadapter "github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/secondary/dynamodb"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	pkglogger "github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

// ── mocks ────────────────────────────────────────────────────────────────────

type mockDynamoDB struct{ mock.Mock }

func (m *mockDynamoDB) PutItem(ctx context.Context, params *awsdynamodb.PutItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*awsdynamodb.PutItemOutput), args.Error(1)
}
func (m *mockDynamoDB) GetItem(ctx context.Context, params *awsdynamodb.GetItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*awsdynamodb.GetItemOutput), args.Error(1)
}
func (m *mockDynamoDB) Scan(ctx context.Context, params *awsdynamodb.ScanInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.ScanOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*awsdynamodb.ScanOutput), args.Error(1)
}
func (m *mockDynamoDB) UpdateItem(ctx context.Context, params *awsdynamodb.UpdateItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.UpdateItemOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*awsdynamodb.UpdateItemOutput), args.Error(1)
}
func (m *mockDynamoDB) DeleteItem(ctx context.Context, params *awsdynamodb.DeleteItemInput, _ ...func(*awsdynamodb.Options)) (*awsdynamodb.DeleteItemOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*awsdynamodb.DeleteItemOutput), args.Error(1)
}

type mockLogger struct{}

func (l *mockLogger) Info(msg string, args ...any)      {}
func (l *mockLogger) Error(msg string, args ...any)     {}
func (l *mockLogger) Debug(msg string, args ...any)     {}
func (l *mockLogger) Warn(msg string, args ...any)      {}
func (l *mockLogger) With(args ...any) pkglogger.Logger { return l }

func newRepo(client *mockDynamoDB) *dynamodbadapter.Repository {
	return dynamodbadapter.New(client, "cauciones-test", &mockLogger{})
}

func sampleCaucion() *domain.Caucion {
	return &domain.Caucion{
		ID:               "test-id",
		Numero:           "CAU-001",
		Tipo:             "garantia",
		Monto:            50_000.0,
		Moneda:           "ARS",
		FechaEmision:     time.Now().UTC(),
		FechaVencimiento: time.Now().UTC().AddDate(1, 0, 0),
		Estado:           domain.EstadoPendiente,
		Beneficiario:     "Empresa SA",
		Tomador:          "Cliente SRL",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
}

// ── Save ─────────────────────────────────────────────────────────────────────

func TestSave_Success(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("PutItem", mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput")).
		Return(&awsdynamodb.PutItemOutput{}, nil)

	err := newRepo(client).Save(context.Background(), sampleCaucion())

	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestSave_ClientError(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("PutItem", mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput")).
		Return(nil, fmt.Errorf("ddb error"))

	err := newRepo(client).Save(context.Background(), sampleCaucion())

	assert.Error(t, err)
}

// ── FindByID ──────────────────────────────────────────────────────────────────

func TestFindByID_Success(t *testing.T) {
	c := sampleCaucion()
	item, _ := attributevalue.MarshalMap(c)
	client := &mockDynamoDB{}
	client.On("GetItem", mock.Anything, mock.AnythingOfType("*dynamodb.GetItemInput")).
		Return(&awsdynamodb.GetItemOutput{Item: item}, nil)

	got, err := newRepo(client).FindByID(context.Background(), "test-id")

	assert.NoError(t, err)
	assert.Equal(t, c.ID, got.ID)
	assert.Equal(t, c.Numero, got.Numero)
}

func TestFindByID_NotFound(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("GetItem", mock.Anything, mock.AnythingOfType("*dynamodb.GetItemInput")).
		Return(&awsdynamodb.GetItemOutput{Item: nil}, nil)

	got, err := newRepo(client).FindByID(context.Background(), "missing")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, got)
}

func TestFindByID_ClientError(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("GetItem", mock.Anything, mock.AnythingOfType("*dynamodb.GetItemInput")).
		Return(nil, fmt.Errorf("ddb error"))

	got, err := newRepo(client).FindByID(context.Background(), "test-id")

	assert.Error(t, err)
	assert.Nil(t, got)
}

// ── FindAll ───────────────────────────────────────────────────────────────────

func TestFindAll_Success(t *testing.T) {
	c := sampleCaucion()
	item, _ := attributevalue.MarshalMap(c)
	client := &mockDynamoDB{}
	client.On("Scan", mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput")).
		Return(&awsdynamodb.ScanOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	cauciones, err := newRepo(client).FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, cauciones, 1)
}

func TestFindAll_Empty(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("Scan", mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput")).
		Return(&awsdynamodb.ScanOutput{Items: []map[string]types.AttributeValue{}}, nil)

	cauciones, err := newRepo(client).FindAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, cauciones)
}

func TestFindAll_ClientError(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("Scan", mock.Anything, mock.AnythingOfType("*dynamodb.ScanInput")).
		Return(nil, fmt.Errorf("ddb error"))

	cauciones, err := newRepo(client).FindAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, cauciones)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("PutItem", mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput")).
		Return(&awsdynamodb.PutItemOutput{}, nil)

	err := newRepo(client).Update(context.Background(), sampleCaucion())

	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestUpdate_ClientError(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("PutItem", mock.Anything, mock.AnythingOfType("*dynamodb.PutItemInput")).
		Return(nil, fmt.Errorf("ddb error"))

	err := newRepo(client).Update(context.Background(), sampleCaucion())

	assert.Error(t, err)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("DeleteItem", mock.Anything, mock.AnythingOfType("*dynamodb.DeleteItemInput")).
		Return(&awsdynamodb.DeleteItemOutput{}, nil)

	err := newRepo(client).Delete(context.Background(), "test-id")

	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestDelete_ClientError(t *testing.T) {
	client := &mockDynamoDB{}
	client.On("DeleteItem", mock.Anything, mock.AnythingOfType("*dynamodb.DeleteItemInput")).
		Return(nil, fmt.Errorf("ddb error"))

	err := newRepo(client).Delete(context.Background(), "test-id")

	assert.Error(t, err)
}
