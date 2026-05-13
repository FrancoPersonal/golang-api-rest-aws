package application_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/application"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/primary"
	pkglogger "github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

// ── mocks ────────────────────────────────────────────────────────────────────

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Save(ctx context.Context, c *domain.Caucion) error {
	return m.Called(ctx, c).Error(0)
}
func (m *mockRepo) FindByID(ctx context.Context, id string) (*domain.Caucion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Caucion), args.Error(1)
}
func (m *mockRepo) FindAll(ctx context.Context) ([]domain.Caucion, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Caucion), args.Error(1)
}
func (m *mockRepo) Update(ctx context.Context, c *domain.Caucion) error {
	return m.Called(ctx, c).Error(0)
}
func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockLogger struct{}

func (l *mockLogger) Info(msg string, args ...any)      {}
func (l *mockLogger) Error(msg string, args ...any)     {}
func (l *mockLogger) Debug(msg string, args ...any)     {}
func (l *mockLogger) Warn(msg string, args ...any)      {}
func (l *mockLogger) With(args ...any) pkglogger.Logger { return l }

func newUC(repo *mockRepo) primary.CaucionUseCase {
	return application.New(repo, &mockLogger{})
}

func sampleInput() primary.CreateCaucionInput {
	return primary.CreateCaucionInput{
		Numero:           "CAU-001",
		Tipo:             "garantia",
		Monto:            100_000.00,
		Moneda:           "ARS",
		FechaEmision:     time.Now().UTC(),
		FechaVencimiento: time.Now().UTC().AddDate(1, 0, 0),
		Beneficiario:     "Empresa SA",
		Tomador:          "Cliente SRL",
	}
}

func sampleCaucion(id string) *domain.Caucion {
	return &domain.Caucion{
		ID:               id,
		Numero:           "CAU-001",
		Tipo:             "garantia",
		Monto:            100_000.00,
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

// ── Create ───────────────────────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	repo := &mockRepo{}
	repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Caucion")).Return(nil)

	c, err := newUC(repo).Create(context.Background(), sampleInput())

	assert.NoError(t, err)
	assert.NotEmpty(t, c.ID)
	assert.Equal(t, domain.EstadoPendiente, c.Estado)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError(t *testing.T) {
	repo := &mockRepo{}
	repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Caucion")).Return(fmt.Errorf("db error"))

	c, err := newUC(repo).Create(context.Background(), sampleInput())

	assert.Error(t, err)
	assert.Nil(t, c)
	repo.AssertExpectations(t)
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	repo := &mockRepo{}
	expected := sampleCaucion("id-1")
	repo.On("FindByID", mock.Anything, "id-1").Return(expected, nil)

	c, err := newUC(repo).GetByID(context.Background(), "id-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, c)
	repo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{}
	repo.On("FindByID", mock.Anything, "missing").Return(nil, domain.ErrNotFound)

	c, err := newUC(repo).GetByID(context.Background(), "missing")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, c)
}

// ── List ─────────────────────────────────────────────────────────────────────

func TestList_Success(t *testing.T) {
	repo := &mockRepo{}
	expected := []domain.Caucion{*sampleCaucion("id-1"), *sampleCaucion("id-2")}
	repo.On("FindAll", mock.Anything).Return(expected, nil)

	cauciones, err := newUC(repo).List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, cauciones, 2)
	repo.AssertExpectations(t)
}

func TestList_Empty(t *testing.T) {
	repo := &mockRepo{}
	repo.On("FindAll", mock.Anything).Return([]domain.Caucion{}, nil)

	cauciones, err := newUC(repo).List(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, cauciones)
}

func TestList_RepoError(t *testing.T) {
	repo := &mockRepo{}
	repo.On("FindAll", mock.Anything).Return(nil, fmt.Errorf("db error"))

	cauciones, err := newUC(repo).List(context.Background())

	assert.Error(t, err)
	assert.Nil(t, cauciones)
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	repo := &mockRepo{}
	existing := sampleCaucion("id-1")
	repo.On("FindByID", mock.Anything, "id-1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Caucion")).Return(nil)

	newNumero := "CAU-002"
	c, err := newUC(repo).Update(context.Background(), "id-1", primary.UpdateCaucionInput{Numero: &newNumero})

	assert.NoError(t, err)
	assert.Equal(t, "CAU-002", c.Numero)
	repo.AssertExpectations(t)
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockRepo{}
	repo.On("FindByID", mock.Anything, "missing").Return(nil, domain.ErrNotFound)

	c, err := newUC(repo).Update(context.Background(), "missing", primary.UpdateCaucionInput{})

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, c)
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	repo := &mockRepo{}
	existing := sampleCaucion("id-1")
	repo.On("FindByID", mock.Anything, "id-1").Return(existing, nil)
	repo.On("Delete", mock.Anything, "id-1").Return(nil)

	err := newUC(repo).Delete(context.Background(), "id-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepo{}
	repo.On("FindByID", mock.Anything, "missing").Return(nil, domain.ErrNotFound)

	err := newUC(repo).Delete(context.Background(), "missing")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// ── ChangeState ───────────────────────────────────────────────────────────────

func TestChangeState_Success(t *testing.T) {
	repo := &mockRepo{}
	existing := sampleCaucion("id-1") // estado: pendiente
	repo.On("FindByID", mock.Anything, "id-1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Caucion")).Return(nil)

	c, err := newUC(repo).ChangeState(context.Background(), "id-1", domain.EstadoVigente)

	assert.NoError(t, err)
	assert.Equal(t, domain.EstadoVigente, c.Estado)
	repo.AssertExpectations(t)
}

func TestChangeState_InvalidTransition(t *testing.T) {
	repo := &mockRepo{}
	existing := sampleCaucion("id-1") // estado: pendiente
	repo.On("FindByID", mock.Anything, "id-1").Return(existing, nil)

	c, err := newUC(repo).ChangeState(context.Background(), "id-1", domain.EstadoVencida)

	assert.ErrorIs(t, err, domain.ErrInvalidTransition)
	assert.Nil(t, c)
}

func TestChangeState_InvalidState(t *testing.T) {
	repo := &mockRepo{}

	c, err := newUC(repo).ChangeState(context.Background(), "id-1", domain.Estado("desconocido"))

	assert.ErrorIs(t, err, domain.ErrInvalidState)
	assert.Nil(t, c)
}

func TestChangeState_NotFound(t *testing.T) {
	repo := &mockRepo{}
	repo.On("FindByID", mock.Anything, "missing").Return(nil, domain.ErrNotFound)

	c, err := newUC(repo).ChangeState(context.Background(), "missing", domain.EstadoVigente)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Nil(t, c)
}

// ── Domain: Estado transitions ────────────────────────────────────────────────

func TestEstado_CanTransitionTo(t *testing.T) {
	cases := []struct {
		from   domain.Estado
		to     domain.Estado
		expect bool
	}{
		{domain.EstadoPendiente, domain.EstadoVigente, true},
		{domain.EstadoPendiente, domain.EstadoCancelada, true},
		{domain.EstadoPendiente, domain.EstadoVencida, false},
		{domain.EstadoVigente, domain.EstadoVencida, true},
		{domain.EstadoVigente, domain.EstadoCancelada, true},
		{domain.EstadoVigente, domain.EstadoPendiente, false},
		{domain.EstadoVencida, domain.EstadoCancelada, false},
	}
	for _, tc := range cases {
		got := tc.from.CanTransitionTo(tc.to)
		assert.Equal(t, tc.expect, got, "%s -> %s", tc.from, tc.to)
	}
}

// ── Domain: sentinel errors ───────────────────────────────────────────────────

func TestDomainErrors_AreDistinct(t *testing.T) {
	assert.False(t, errors.Is(domain.ErrNotFound, domain.ErrInvalidState))
	assert.False(t, errors.Is(domain.ErrInvalidState, domain.ErrInvalidTransition))
}
