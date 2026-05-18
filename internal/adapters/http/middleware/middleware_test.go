package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

// buildJWT constructs a minimal HS256 JWT for use in tests.
func buildJWT(t *testing.T, secret string, exp int64) string {
	t.Helper()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadData, err := json.Marshal(map[string]any{"exp": exp})
	require.NoError(t, err)

	payload := base64.RawURLEncoding.EncodeToString(payloadData)
	signingInput := header + "." + payload

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig
}

// noopHandler is a Handler that returns 200 OK.
var noopHandler Handler = func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

// errHandler is a Handler that returns a 500 and an error.
var errHandler Handler = func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("boom")
}

// ─── Chain ────────────────────────────────────────────────────────────────────

func TestChainNoMiddlewares(t *testing.T) {
	h := Chain(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestChainSingleMiddleware(t *testing.T) {
	var called bool
	mw := func(next Handler) Handler {
		return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			called = true
			return next(ctx, req)
		}
	}
	h := Chain(noopHandler, mw)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.True(t, called)
}

func TestChainMultipleMiddlewaresOrder(t *testing.T) {
	var order []int
	mw := func(n int) func(Handler) Handler {
		return func(next Handler) Handler {
			return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
				order = append(order, n)
				return next(ctx, req)
			}
		}
	}
	h := Chain(noopHandler, mw(1), mw(2))
	_, _ = h(context.Background(), events.APIGatewayProxyRequest{})
	require.Equal(t, []int{1, 2}, order)
}

// ─── LoggingMiddleware ────────────────────────────────────────────────────────

func TestLoggingMiddlewareSuccess(t *testing.T) {
	var buf strings.Builder
	logger := log.New(&buf, "", 0)

	h := LoggingMiddleware(logger)(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/test",
	})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Contains(t, buf.String(), "level=info")
}

func TestLoggingMiddlewareWithError(t *testing.T) {
	var buf strings.Builder
	logger := log.New(&buf, "", 0)

	h := LoggingMiddleware(logger)(errHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/test",
	})
	require.Error(t, err)
	require.Equal(t, 500, resp.StatusCode)
	require.Contains(t, buf.String(), "level=error")
}

func TestLoggingMiddlewareNilLogger(t *testing.T) {
	h := LoggingMiddleware(nil)(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

// ─── JWTAuthMiddleware ────────────────────────────────────────────────────────

func TestJWTAuthMiddlewareMissingHeader(t *testing.T) {
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuthMiddlewareInvalidBearerFormat(t *testing.T) {
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": "InvalidFormat"},
	})
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuthMiddlewareInvalidToken(t *testing.T) {
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": "Bearer not.a.valid.jwt"},
	})
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuthMiddlewareWrongSecret(t *testing.T) {
	token := buildJWT(t, "other-secret", time.Now().Add(5*time.Minute).Unix())
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": "Bearer " + token},
	})
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuthMiddlewareExpiredToken(t *testing.T) {
	token := buildJWT(t, "secret", time.Now().Add(-5*time.Minute).Unix())
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": "Bearer " + token},
	})
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

func TestJWTAuthMiddlewareValidToken(t *testing.T) {
	token := buildJWT(t, "secret", time.Now().Add(5*time.Minute).Unix())
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": "Bearer " + token},
	})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestJWTAuthMiddlewareLowercaseAuthHeader(t *testing.T) {
	token := buildJWT(t, "secret", time.Now().Add(5*time.Minute).Unix())
	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"authorization": "Bearer " + token},
	})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestJWTAuthMiddlewareTokenWithoutExp(t *testing.T) {
	// Build a JWT with no exp claim.
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, _ := json.Marshal(map[string]any{})
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := header + "." + payload
	mac := hmac.New(sha256.New, []byte("secret"))
	_, _ = mac.Write([]byte(signingInput))
	token := signingInput + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	h := JWTAuthMiddleware("secret")(noopHandler)
	resp, err := h(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": "Bearer " + token},
	})
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
