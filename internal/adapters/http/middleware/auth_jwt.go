package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/dto"
)

type jwtClaims struct {
	Exp int64 `json:"exp"`
}

func JWTAuthMiddleware(secret string) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			authHeader := req.Headers["Authorization"]
			if authHeader == "" {
				authHeader = req.Headers["authorization"]
			}

			token, ok := extractBearerToken(authHeader)
			if !ok {
				return dto.Fail(401, "unauthorized", "missing or invalid authorization header"), nil
			}

			if err := validateJWT(token, secret); err != nil {
				return dto.Fail(401, "unauthorized", "invalid jwt token"), nil
			}

			return next(ctx, req)
		}
	}
}

func extractBearerToken(header string) (string, bool) {
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}
	return strings.TrimSpace(parts[1]), true
}

func validateJWT(token, secret string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid token format")
	}

	if secret == "" {
		return errors.New("missing jwt secret")
	}

	signingInput := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	expected := mac.Sum(nil)
	if !hmac.Equal(sig, expected) {
		return errors.New("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	var claims jwtClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return err
	}

	if claims.Exp > 0 && time.Now().Unix() >= claims.Exp {
		return errors.New("token expired")
	}

	return nil
}
