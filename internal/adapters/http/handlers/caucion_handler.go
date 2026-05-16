package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/dto"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	usecasep "github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/usecase"
)

type CaucionHandler struct {
	useCase usecasep.CaucionUseCase
}

func NewCaucionHandler(useCase usecasep.CaucionUseCase) *CaucionHandler {
	return &CaucionHandler{useCase: useCase}
}

func (h *CaucionHandler) CreateCaucion(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var in dto.CreateCaucionRequest
	if err := json.Unmarshal([]byte(req.Body), &in); err != nil {
		return dto.Fail(400, "invalid_json", "request body must be valid JSON"), nil
	}

	created, err := h.useCase.Create(ctx, domain.CreateCaucionInput{
		Numero:           in.Numero,
		Tipo:             in.Tipo,
		Monto:            in.Monto,
		Moneda:           in.Moneda,
		Beneficiario:     in.Beneficiario,
		Tomador:          in.Tomador,
		FechaEmision:     in.FechaEmision,
		FechaVencimiento: in.FechaVencimiento,
	})
	if err != nil {
		return mapError(err), nil
	}

	resp := dto.Success(201, created)
	if resp.Headers == nil {
		resp.Headers = map[string]string{}
	}
	resp.Headers["Location"] = "/cauciones/" + created.ID

	return resp, nil
}

func (h *CaucionHandler) GetCauciones(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	items, err := h.useCase.List(ctx)
	if err != nil {
		return mapError(err), nil
	}

	return dto.Success(200, items), nil
}

func mapError(err error) events.APIGatewayProxyResponse {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return dto.Fail(400, "invalid_input", err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return dto.Fail(401, "unauthorized", "unauthorized")
	case errors.Is(err, domain.ErrNotFound):
		return dto.Fail(404, "not_found", "resource not found")
	case errors.Is(err, domain.ErrConflict):
		return dto.Fail(409, "conflict", "resource conflict")
	default:
		return dto.Fail(500, "internal_error", "internal server error")
	}
}
