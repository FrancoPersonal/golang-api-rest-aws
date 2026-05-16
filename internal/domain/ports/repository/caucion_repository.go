package repository

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type CaucionRepository interface {
	Create(ctx context.Context, caucion domain.Caucion) error
	List(ctx context.Context) ([]domain.Caucion, error)
}
