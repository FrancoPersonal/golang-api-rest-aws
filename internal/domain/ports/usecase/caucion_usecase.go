package usecase

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type SuretyBondUseCase interface {
	Create(ctx context.Context, input domain.CreateSuretyBondInput) (domain.SuretyBond, error)
	List(ctx context.Context) ([]domain.SuretyBond, error)
}
