package sqlserver

import (
	"context"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type SuretyBondRepository struct{}

func NewSuretyBondRepository() *SuretyBondRepository {
	return &SuretyBondRepository{}
}

func (r *SuretyBondRepository) Create(_ context.Context, _ domain.SuretyBond) error {
	return domain.ErrInternal
}
