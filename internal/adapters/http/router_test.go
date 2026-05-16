package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/handlers"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
)

type routerUseCaseMock struct{}

func (routerUseCaseMock) Create(_ context.Context, _ domain.CreateCaucionInput) (domain.Caucion, error) {
	return domain.Caucion{ID: "id-1", Numero: "C-001", CreatedAt: time.Now().UTC()}, nil
}

func (routerUseCaseMock) List(_ context.Context) ([]domain.Caucion, error) {
	return []domain.Caucion{{ID: "id-1"}}, nil
}

func TestRouterJWTMiddleware(t *testing.T) {
	h := handlers.NewCaucionHandler(routerUseCaseMock{})
	r := NewRouter(h, nil, "secret")

	body := `{"numero":"C-1","tipo":"tradicional","monto":1000,"moneda":"ARS","beneficiario":"Banco","tomador":"Cliente","fecha_emision":"2026-01-01T00:00:00Z","fecha_vencimiento":"2026-01-02T00:00:00Z"}`

	unauthorizedResp, err := r.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/cauciones",
		Body:       body,
	})
	require.NoError(t, err)
	require.Equal(t, 401, unauthorizedResp.StatusCode)

	token := signJWT(t, "secret", time.Now().Add(5*time.Minute).Unix())
	authorizedResp, err := r.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/cauciones",
		Body:       body,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 201, authorizedResp.StatusCode)

	unauthorizedGetResp, err := r.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/cauciones",
	})
	require.NoError(t, err)
	require.Equal(t, 401, unauthorizedGetResp.StatusCode)

	authorizedGetResp, err := r.Handle(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/cauciones",
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 200, authorizedGetResp.StatusCode)
}

func signJWT(t *testing.T, secret string, exp int64) string {
	t.Helper()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadData, err := json.Marshal(map[string]any{"exp": exp})
	require.NoError(t, err)

	payload := base64.RawURLEncoding.EncodeToString(payloadData)
	signingInput := header + "." + payload

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature
}
