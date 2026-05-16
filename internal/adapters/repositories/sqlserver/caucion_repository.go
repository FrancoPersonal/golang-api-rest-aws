package sqlserver

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type CaucionRepository struct{}

func NewCaucionRepository() *CaucionRepository {
	return &CaucionRepository{}
}

func (r *CaucionRepository) Create(_ context.Context, _ domain.Caucion) error {
	return domain.ErrInternal
}
