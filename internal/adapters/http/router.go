package http

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/dto"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/handlers"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/http/middleware"
)

type Router struct {
	createCaucion middleware.Handler
	listCauciones middleware.Handler
}

func NewRouter(handler *handlers.CaucionHandler, logger *log.Logger, jwtSecret string) *Router {
	create := middleware.Chain(
		handler.CreateCaucion,
		middleware.LoggingMiddleware(logger),
		middleware.JWTAuthMiddleware(jwtSecret),
	)

	list := middleware.Chain(
		handler.GetCauciones,
		middleware.LoggingMiddleware(logger),
		middleware.JWTAuthMiddleware(jwtSecret),
	)

	return &Router{createCaucion: create, listCauciones: list}
}

func (r *Router) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToUpper(req.HTTPMethod)
	path := strings.TrimSuffix(req.Path, "/")
	if path == "" {
		path = "/"
	}

	if method == "POST" && path == "/cauciones" {
		return r.createCaucion(ctx, req)
	}

	if method == "GET" && path == "/cauciones" {
		return r.listCauciones(ctx, req)
	}

	return dto.Fail(404, "not_found", "route not found"), nil
}
