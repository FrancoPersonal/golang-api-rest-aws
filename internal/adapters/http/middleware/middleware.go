package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Handler func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func Chain(base Handler, middlewares ...func(Handler) Handler) Handler {
	wrapped := base
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
