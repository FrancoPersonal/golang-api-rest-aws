package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type repoMock struct {
	createErr error
	created   domain.Caucion
	list      []domain.Caucion
	listErr   error
}

func (m *repoMock) Create(_ context.Context, caucion domain.Caucion) error {
	m.created = caucion
	return m.createErr
}

func (m *repoMock) List(_ context.Context) ([]domain.Caucion, error) {
	return m.list, m.listErr
}

func TestCaucionServiceCreateSuccess(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	repo := &repoMock{}
	svc := NewCaucionService(repo, func() time.Time { return now }, func() string { return "id-123" })

	result, err := svc.Create(context.Background(), domain.CreateCaucionInput{
		Numero:           "C-001",
		Tipo:             "tradicional",
		Monto:            1500,
		Moneda:           "ars",
		Beneficiario:     "Banco",
		Tomador:          "Cliente",
		FechaEmision:     now,
		FechaVencimiento: now.Add(24 * time.Hour),
	})

	require.NoError(t, err)
	require.Equal(t, "id-123", result.ID)
	require.Equal(t, domain.EstadoPendiente, result.Estado)
	require.Equal(t, "ARS", result.Moneda)
	require.Equal(t, repo.created.ID, result.ID)
}

func TestCaucionServiceCreateInvalidInput(t *testing.T) {
	repo := &repoMock{}
	svc := NewCaucionService(repo, nil, nil)

	_, err := svc.Create(context.Background(), domain.CreateCaucionInput{})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCaucionServiceCreateRepoError(t *testing.T) {
	repo := &repoMock{createErr: domain.ErrConflict}
	svc := NewCaucionService(repo, nil, nil)

	_, err := svc.Create(context.Background(), domain.CreateCaucionInput{
		Numero:           "C-001",
		Tipo:             "tradicional",
		Monto:            1500,
		Moneda:           "ARS",
		Beneficiario:     "Banco",
		Tomador:          "Cliente",
		FechaEmision:     time.Now().UTC(),
		FechaVencimiento: time.Now().UTC().Add(24 * time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrConflict)
}

func TestCaucionServiceListSuccess(t *testing.T) {
	repo := &repoMock{list: []domain.Caucion{{ID: "1"}, {ID: "2"}}}
	svc := NewCaucionService(repo, nil, nil)

	items, err := svc.List(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 2)
}

func TestCaucionServiceListError(t *testing.T) {
	repo := &repoMock{listErr: domain.ErrInternal}
	svc := NewCaucionService(repo, nil, nil)

	items, err := svc.List(context.Background())
	require.Nil(t, items)
	require.ErrorIs(t, err, domain.ErrInternal)
}
