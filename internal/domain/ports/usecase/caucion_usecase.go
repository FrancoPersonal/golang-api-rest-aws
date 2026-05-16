package usecase

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type CaucionUseCase interface {
	Create(ctx context.Context, input domain.CreateCaucionInput) (domain.Caucion, error)
	List(ctx context.Context) ([]domain.Caucion, error)
}
