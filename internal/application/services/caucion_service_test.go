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
	created   domain.SuretyBond
	list      []domain.SuretyBond
	listErr   error
}

func (m *repoMock) Create(_ context.Context, suretyBond domain.SuretyBond) error {
	m.created = suretyBond
	return m.createErr
}

func (m *repoMock) List(_ context.Context) ([]domain.SuretyBond, error) {
	return m.list, m.listErr
}

func TestSuretyBondServiceCreateSuccess(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	repo := &repoMock{}
	svc := NewSuretyBondService(repo, func() time.Time { return now }, func() string { return "id-123" })

	result, err := svc.Create(context.Background(), domain.CreateSuretyBondInput{
		Number:      "C-001",
		Type:        "traditional",
		Amount:      1500,
		Currency:    "ars",
		Beneficiary: "Banco",
		Holder:      "Cliente",
		IssueDate:   now,
		ExpiryDate:  now.Add(24 * time.Hour),
	})

	require.NoError(t, err)
	require.Equal(t, "id-123", result.ID)
	require.Equal(t, domain.StatusPending, result.Status)
	require.Equal(t, "ARS", result.Currency)
	require.Equal(t, repo.created.ID, result.ID)
}

func TestSuretyBondServiceCreateInvalidInput(t *testing.T) {
	repo := &repoMock{}
	svc := NewSuretyBondService(repo, nil, nil)

	_, err := svc.Create(context.Background(), domain.CreateSuretyBondInput{})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestSuretyBondServiceCreateRepoError(t *testing.T) {
	repo := &repoMock{createErr: domain.ErrConflict}
	svc := NewSuretyBondService(repo, nil, nil)

	_, err := svc.Create(context.Background(), domain.CreateSuretyBondInput{
		Number:      "C-001",
		Type:        "traditional",
		Amount:      1500,
		Currency:    "ARS",
		Beneficiary: "Banco",
		Holder:      "Cliente",
		IssueDate:   time.Now().UTC(),
		ExpiryDate:  time.Now().UTC().Add(24 * time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrConflict)
}

func TestSuretyBondServiceListSuccess(t *testing.T) {
	repo := &repoMock{list: []domain.SuretyBond{{ID: "1"}, {ID: "2"}}}
	svc := NewSuretyBondService(repo, nil, nil)

	items, err := svc.List(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 2)
}

func TestSuretyBondServiceListError(t *testing.T) {
	repo := &repoMock{listErr: domain.ErrInternal}
	svc := NewSuretyBondService(repo, nil, nil)

	items, err := svc.List(context.Background())
	require.Nil(t, items)
	require.ErrorIs(t, err, domain.ErrInternal)
}

func TestSuretyBondServiceListNilReturnsEmpty(t *testing.T) {
	repo := &repoMock{list: nil}
	svc := NewSuretyBondService(repo, nil, nil)

	items, err := svc.List(context.Background())
	require.NoError(t, err)
	require.NotNil(t, items)
	require.Empty(t, items)
}

func TestValidateCreateInputExpiryBeforeIssue(t *testing.T) {
	repo := &repoMock{}
	svc := NewSuretyBondService(repo, nil, nil)

	now := time.Now().UTC()
	_, err := svc.Create(context.Background(), domain.CreateSuretyBondInput{
		Number:      "C-001",
		Type:        "traditional",
		Amount:      1500,
		Currency:    "ARS",
		Beneficiary: "Banco",
		Holder:      "Cliente",
		IssueDate:   now,
		ExpiryDate:  now.Add(-24 * time.Hour),
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestValidateCreateInputZeroAmount(t *testing.T) {
	repo := &repoMock{}
	svc := NewSuretyBondService(repo, nil, nil)

	_, err := svc.Create(context.Background(), domain.CreateSuretyBondInput{
		Number:      "C-001",
		Type:        "traditional",
		Amount:      0,
		Currency:    "ARS",
		Beneficiary: "Banco",
		Holder:      "Cliente",
		IssueDate:   time.Now().UTC(),
		ExpiryDate:  time.Now().UTC().Add(24 * time.Hour),
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestValidateCreateInputZeroDates(t *testing.T) {
	repo := &repoMock{}
	svc := NewSuretyBondService(repo, nil, nil)

	_, err := svc.Create(context.Background(), domain.CreateSuretyBondInput{
		Number:      "C-001",
		Type:        "traditional",
		Amount:      1500,
		Currency:    "ARS",
		Beneficiary: "Banco",
		Holder:      "Cliente",
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}
