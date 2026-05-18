package middleware

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func LoggingMiddleware(logger *log.Logger) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			duration := time.Since(start)

			if logger != nil {
				if err != nil {
					logger.Printf("level=error method=%s path=%s status=%d duration_ms=%d err=%q", req.HTTPMethod, req.Path, resp.StatusCode, duration.Milliseconds(), err.Error())
				} else {
					logger.Printf("level=info method=%s path=%s status=%d duration_ms=%d", req.HTTPMethod, req.Path, resp.StatusCode, duration.Milliseconds())
				}
			}

			return resp, err
		}
	}
}
