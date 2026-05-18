package repository

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type SuretyBondRepository interface {
	Create(ctx context.Context, suretyBond domain.SuretyBond) error
	List(ctx context.Context) ([]domain.SuretyBond, error)
}
