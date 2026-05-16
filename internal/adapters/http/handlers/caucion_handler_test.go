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
	result domain.Caucion
	err    error
	items  []domain.Caucion
}

func (m *useCaseMock) Create(_ context.Context, _ domain.CreateCaucionInput) (domain.Caucion, error) {
	return m.result, m.err
}

func (m *useCaseMock) List(_ context.Context) ([]domain.Caucion, error) {
	return m.items, m.err
}

func TestCreateCaucionSuccess(t *testing.T) {
	uc := &useCaseMock{result: domain.Caucion{ID: "id-1", Numero: "C-1", CreatedAt: time.Now().UTC()}}
	h := NewCaucionHandler(uc)

	body := `{"numero":"C-1","tipo":"tradicional","monto":1000,"moneda":"ARS","beneficiario":"Banco","tomador":"Cliente","fecha_emision":"2026-01-01T00:00:00Z","fecha_vencimiento":"2026-01-02T00:00:00Z"}`
	resp, err := h.CreateCaucion(context.Background(), events.APIGatewayProxyRequest{Body: body})

	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
	require.Equal(t, "/cauciones/id-1", resp.Headers["Location"])

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	require.Equal(t, true, payload["success"])
}

func TestCreateCaucionInvalidJSON(t *testing.T) {
	uc := &useCaseMock{}
	h := NewCaucionHandler(uc)

	resp, err := h.CreateCaucion(context.Background(), events.APIGatewayProxyRequest{Body: "{"})
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}

func TestGetCaucionesSuccess(t *testing.T) {
	uc := &useCaseMock{items: []domain.Caucion{{ID: "id-1"}, {ID: "id-2"}}}
	h := NewCaucionHandler(uc)

	resp, err := h.GetCauciones(context.Background(), events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	require.Equal(t, true, payload["success"])
}
