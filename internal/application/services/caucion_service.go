package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	repop "github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/repository"
)

type CaucionService struct {
	repo  repop.CaucionRepository
	now   func() time.Time
	newID func() string
}

func NewCaucionService(repo repop.CaucionRepository, now func() time.Time, newID func() string) *CaucionService {
	if now == nil {
		now = time.Now
	}
	if newID == nil {
		newID = uuid.NewString
	}

	return &CaucionService{
		repo:  repo,
		now:   now,
		newID: newID,
	}
}

func (s *CaucionService) Create(ctx context.Context, input domain.CreateCaucionInput) (domain.Caucion, error) {
	if err := validateCreateInput(input); err != nil {
		return domain.Caucion{}, err
	}

	now := s.now().UTC()
	c := domain.Caucion{
		ID:               s.newID(),
		Numero:           strings.TrimSpace(input.Numero),
		Tipo:             strings.TrimSpace(input.Tipo),
		Monto:            input.Monto,
		Moneda:           strings.ToUpper(strings.TrimSpace(input.Moneda)),
		Estado:           domain.EstadoPendiente,
		Beneficiario:     strings.TrimSpace(input.Beneficiario),
		Tomador:          strings.TrimSpace(input.Tomador),
		FechaEmision:     input.FechaEmision.UTC(),
		FechaVencimiento: input.FechaVencimiento.UTC(),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return domain.Caucion{}, err
	}

	return c, nil
}

func (s *CaucionService) List(ctx context.Context) ([]domain.Caucion, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if items == nil {
		return []domain.Caucion{}, nil
	}

	return items, nil
}

func validateCreateInput(input domain.CreateCaucionInput) error {
	if strings.TrimSpace(input.Numero) == "" ||
		strings.TrimSpace(input.Tipo) == "" ||
		strings.TrimSpace(input.Moneda) == "" ||
		strings.TrimSpace(input.Beneficiario) == "" ||
		strings.TrimSpace(input.Tomador) == "" {
		return domain.ErrInvalidInput
	}

	if input.Monto <= 0 {
		return domain.ErrInvalidInput
	}

	if input.FechaEmision.IsZero() || input.FechaVencimiento.IsZero() {
		return domain.ErrInvalidInput
	}

	if input.FechaVencimiento.Before(input.FechaEmision) {
		return domain.ErrInvalidInput
	}

	return nil
}
