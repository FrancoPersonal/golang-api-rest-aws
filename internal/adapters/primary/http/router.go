package httphandler

import (
	"net/http"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/primary"
	"github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

// NewRouter wires all Caucion routes into an http.ServeMux and returns the handler.
// Uses Go 1.22+ enhanced mux with method + path syntax.
func NewRouter(uc primary.CaucionUseCase, log logger.Logger) http.Handler {
	h := &Handler{uc: uc, log: log}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /cauciones", h.Create)
	mux.HandleFunc("GET /cauciones", h.List)
	mux.HandleFunc("GET /cauciones/{id}", h.GetByID)
	mux.HandleFunc("PUT /cauciones/{id}", h.Update)
	mux.HandleFunc("DELETE /cauciones/{id}", h.Delete)
	mux.HandleFunc("PATCH /cauciones/{id}/estado", h.ChangeState)

	return mux
}
