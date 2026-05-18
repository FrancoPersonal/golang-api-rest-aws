package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	repop "github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/repository"
)

type SuretyBondService struct {
	repo  repop.SuretyBondRepository
	now   func() time.Time
	newID func() string
}

func NewSuretyBondService(repo repop.SuretyBondRepository, now func() time.Time, newID func() string) *SuretyBondService {
	if now == nil {
		now = time.Now
	}
	if newID == nil {
		newID = uuid.NewString
	}

	return &SuretyBondService{
		repo:  repo,
		now:   now,
		newID: newID,
	}
}

func (s *SuretyBondService) Create(ctx context.Context, input domain.CreateSuretyBondInput) (domain.SuretyBond, error) {
	if err := validateCreateInput(input); err != nil {
		return domain.SuretyBond{}, err
	}

	now := s.now().UTC()
	c := domain.SuretyBond{
		ID:          s.newID(),
		Number:      strings.TrimSpace(input.Number),
		Type:        strings.TrimSpace(input.Type),
		Amount:      input.Amount,
		Currency:    strings.ToUpper(strings.TrimSpace(input.Currency)),
		Status:      domain.StatusPending,
		Beneficiary: strings.TrimSpace(input.Beneficiary),
		Holder:      strings.TrimSpace(input.Holder),
		IssueDate:   input.IssueDate.UTC(),
		ExpiryDate:  input.ExpiryDate.UTC(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return domain.SuretyBond{}, err
	}

	return c, nil
}

func (s *SuretyBondService) List(ctx context.Context) ([]domain.SuretyBond, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if items == nil {
		return []domain.SuretyBond{}, nil
	}

	return items, nil
}

func validateCreateInput(input domain.CreateSuretyBondInput) error {
	if strings.TrimSpace(input.Number) == "" ||
		strings.TrimSpace(input.Type) == "" ||
		strings.TrimSpace(input.Currency) == "" ||
		strings.TrimSpace(input.Beneficiary) == "" ||
		strings.TrimSpace(input.Holder) == "" {
		return domain.ErrInvalidInput
	}

	if input.Amount <= 0 {
		return domain.ErrInvalidInput
	}

	if input.IssueDate.IsZero() || input.ExpiryDate.IsZero() {
		return domain.ErrInvalidInput
	}

	if input.ExpiryDate.Before(input.IssueDate) {
		return domain.ErrInvalidInput
	}

	return nil
}
