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
	createSuretyBond middleware.Handler
	listSuretyBonds  middleware.Handler
}

func NewRouter(handler *handlers.SuretyBondHandler, logger *log.Logger, jwtSecret string) *Router {
	create := middleware.Chain(
		handler.CreateSuretyBond,
		middleware.LoggingMiddleware(logger),
		middleware.JWTAuthMiddleware(jwtSecret),
	)

	list := middleware.Chain(
		handler.ListSuretyBonds,
		middleware.LoggingMiddleware(logger),
		middleware.JWTAuthMiddleware(jwtSecret),
	)

	return &Router{createSuretyBond: create, listSuretyBonds: list}
}

func (r *Router) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := strings.ToUpper(req.HTTPMethod)
	path := strings.TrimSuffix(req.Path, "/")
	if path == "" {
		path = "/"
	}

	if method == "POST" && path == "/suretyBonds" {
		return r.createSuretyBond(ctx, req)
	}

	if method == "GET" && path == "/suretyBonds" {
		return r.listSuretyBonds(ctx, req)
	}

	return dto.Fail(404, "not_found", "route not found"), nil
}
