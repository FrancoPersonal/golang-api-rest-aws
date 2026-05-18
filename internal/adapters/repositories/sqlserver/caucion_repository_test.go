package sqlserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

func TestCreateReturnsInternalError(t *testing.T) {
	repo := NewSuretyBondRepository()
	err := repo.Create(context.Background(), domain.SuretyBond{ID: "id"})
	require.ErrorIs(t, err, domain.ErrInternal)
}
