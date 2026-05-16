package handlers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type useCaseMock struct {
	result domain.SuretyBond
	err    error
	items  []domain.SuretyBond
}

func (m *useCaseMock) Create(_ context.Context, _ domain.CreateSuretyBondInput) (domain.SuretyBond, error) {
	return m.result, m.err
}

func (m *useCaseMock) List(_ context.Context) ([]domain.SuretyBond, error) {
	return m.items, m.err
}

func TestCreateSuretyBondSuccess(t *testing.T) {
	uc := &useCaseMock{result: domain.SuretyBond{ID: "id-1", Number: "C-1", CreatedAt: time.Now().UTC()}}
	h := NewSuretyBondHandler(uc)

	body := `{"numero":"C-1","tipo":"traditional","monto":1000,"moneda":"ARS","beneficiario":"Banco","tomador":"Cliente","fecha_emision":"2026-01-01T00:00:00Z","fecha_vencimiento":"2026-01-02T00:00:00Z"}`
	resp, err := h.CreateSuretyBond(context.Background(), events.APIGatewayProxyRequest{Body: body})

	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
	require.Equal(t, "/suretyBonds/id-1", resp.Headers["Location"])

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	require.Equal(t, true, payload["success"])
}

func TestCreateSuretyBondInvalidJSON(t *testing.T) {
	uc := &useCaseMock{}
	h := NewSuretyBondHandler(uc)

	resp, err := h.CreateSuretyBond(context.Background(), events.APIGatewayProxyRequest{Body: "{"})
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

func TestListSuretyBondsSuccess(t *testing.T) {
	uc := &useCaseMock{items: []domain.SuretyBond{{ID: "id-1"}, {ID: "id-2"}}}
	h := NewSuretyBondHandler(uc)

	resp, err := h.ListSuretyBonds(context.Background(), events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	require.Equal(t, true, payload["success"])
}
