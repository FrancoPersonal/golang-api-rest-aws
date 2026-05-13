package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/primary"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/secondary"
	"github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

type caucionUseCase struct {
	repo secondary.CaucionRepository
	log  logger.Logger
}

// New creates a new CaucionUseCase that implements the primary.CaucionUseCase port.
func New(repo secondary.CaucionRepository, log logger.Logger) primary.CaucionUseCase {
	return &caucionUseCase{repo: repo, log: log}
}

func (uc *caucionUseCase) Create(ctx context.Context, input primary.CreateCaucionInput) (*domain.Caucion, error) {
	now := time.Now().UTC()
	c := &domain.Caucion{
		ID:               uuid.NewString(),
		Numero:           input.Numero,
		Tipo:             input.Tipo,
		Monto:            input.Monto,
		Moneda:           input.Moneda,
		FechaEmision:     input.FechaEmision,
		FechaVencimiento: input.FechaVencimiento,
		Estado:           domain.EstadoPendiente,
		Beneficiario:     input.Beneficiario,
		Tomador:          input.Tomador,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := uc.repo.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("create caucion: %w", err)
	}

	uc.log.Info("caucion created", "id", c.ID)
	return c, nil
}

func (uc *caucionUseCase) GetByID(ctx context.Context, id string) (*domain.Caucion, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get caucion %s: %w", id, err)
	}
	return c, nil
}

func (uc *caucionUseCase) List(ctx context.Context) ([]domain.Caucion, error) {
	cauciones, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cauciones: %w", err)
	}
	return cauciones, nil
}

func (uc *caucionUseCase) Update(ctx context.Context, id string, input primary.UpdateCaucionInput) (*domain.Caucion, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("update caucion %s: %w", id, err)
	}

	if input.Numero != nil {
		c.Numero = *input.Numero
	}
	if input.Tipo != nil {
		c.Tipo = *input.Tipo
	}
	if input.Monto != nil {
		c.Monto = *input.Monto
	}
	if input.Moneda != nil {
		c.Moneda = *input.Moneda
	}
	if input.FechaEmision != nil {
		c.FechaEmision = *input.FechaEmision
	}
	if input.FechaVencimiento != nil {
		c.FechaVencimiento = *input.FechaVencimiento
	}
	if input.Beneficiario != nil {
		c.Beneficiario = *input.Beneficiario
	}
	if input.Tomador != nil {
		c.Tomador = *input.Tomador
	}
	c.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("update caucion %s: %w", id, err)
	}

	uc.log.Info("caucion updated", "id", id)
	return c, nil
}

func (uc *caucionUseCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("delete caucion %s: %w", id, err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete caucion %s: %w", id, err)
	}

	uc.log.Info("caucion deleted", "id", id)
	return nil
}

func (uc *caucionUseCase) ChangeState(ctx context.Context, id string, newState domain.Estado) (*domain.Caucion, error) {
	switch newState {
	case domain.EstadoPendiente, domain.EstadoVigente, domain.EstadoVencida, domain.EstadoCancelada:
		// valid estado
	default:
		return nil, domain.ErrInvalidState
	}

	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("change state caucion %s: %w", id, err)
	}

	if !c.Estado.CanTransitionTo(newState) {
		return nil, fmt.Errorf("change state caucion %s from %s to %s: %w", id, c.Estado, newState, domain.ErrInvalidTransition)
	}

	c.Estado = newState
	c.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("change state caucion %s: %w", id, err)
	}

	uc.log.Info("caucion state changed", "id", id, "estado", newState)
	return c, nil
}
