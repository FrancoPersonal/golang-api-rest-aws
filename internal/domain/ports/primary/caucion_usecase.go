package primary

import (
	"context"
	"time"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

// CreateCaucionInput holds the data required to create a new Caucion.
type CreateCaucionInput struct {
	Numero           string
	Tipo             string
	Monto            float64
	Moneda           string
	FechaEmision     time.Time
	FechaVencimiento time.Time
	Beneficiario     string
	Tomador          string
}

// UpdateCaucionInput holds the updatable fields of a Caucion (all optional via pointers).
type UpdateCaucionInput struct {
	Numero           *string
	Tipo             *string
	Monto            *float64
	Moneda           *string
	FechaEmision     *time.Time
	FechaVencimiento *time.Time
	Beneficiario     *string
	Tomador          *string
}

// CaucionUseCase is the primary (input) port for all Caucion operations.
type CaucionUseCase interface {
	Create(ctx context.Context, input CreateCaucionInput) (*domain.Caucion, error)
	GetByID(ctx context.Context, id string) (*domain.Caucion, error)
	List(ctx context.Context) ([]domain.Caucion, error)
	Update(ctx context.Context, id string, input UpdateCaucionInput) (*domain.Caucion, error)
	Delete(ctx context.Context, id string) error
	ChangeState(ctx context.Context, id string, newState domain.Estado) (*domain.Caucion, error)
}
