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

type SuretyBondHandler struct {
	useCase usecasep.SuretyBondUseCase
}

func NewSuretyBondHandler(useCase usecasep.SuretyBondUseCase) *SuretyBondHandler {
	return &SuretyBondHandler{useCase: useCase}
}

func (h *SuretyBondHandler) CreateSuretyBond(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var in dto.CreateSuretyBondRequest
	if err := json.Unmarshal([]byte(req.Body), &in); err != nil {
		return dto.Fail(400, "invalid_json", "request body must be valid JSON"), nil
	}

	created, err := h.useCase.Create(ctx, domain.CreateSuretyBondInput{
		Number:      in.Number,
		Type:        in.Type,
		Amount:      in.Amount,
		Currency:    in.Currency,
		Beneficiary: in.Beneficiary,
		Holder:      in.Holder,
		IssueDate:   in.IssueDate,
		ExpiryDate:  in.ExpiryDate,
	})
	if err != nil {
		return mapError(err), nil
	}

	resp := dto.Success(201, created)
	if resp.Headers == nil {
		resp.Headers = map[string]string{}
	}
	resp.Headers["Location"] = "/suretyBonds/" + created.ID

	return resp, nil
}

func (h *SuretyBondHandler) ListSuretyBonds(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
