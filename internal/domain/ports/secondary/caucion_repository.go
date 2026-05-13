package secondary

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

// CaucionRepository is the secondary (output) port for persistence operations.
type CaucionRepository interface {
	Save(ctx context.Context, caucion *domain.Caucion) error
	FindByID(ctx context.Context, id string) (*domain.Caucion, error)
	FindAll(ctx context.Context) ([]domain.Caucion, error)
	Update(ctx context.Context, caucion *domain.Caucion) error
	Delete(ctx context.Context, id string) error
}
