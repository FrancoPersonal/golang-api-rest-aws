package httphandler

import (
	"errors"
	"net/http"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/primary"
	"github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

// Handler holds the HTTP handlers for the Cauciones resource.
type Handler struct {
	uc  primary.CaucionUseCase
	log logger.Logger
}

// Create handles POST /cauciones.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := decodeJSON[CreateCaucionRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := req.validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	c, err := h.uc.Create(r.Context(), req.toInput())
	if err != nil {
		h.log.Error("handler: create caucion", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

// GetByID handles GET /cauciones/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	c, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "caucion not found")
			return
		}
		h.log.Error("handler: get caucion", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, c)
}

// List handles GET /cauciones.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cauciones, err := h.uc.List(r.Context())
	if err != nil {
		h.log.Error("handler: list cauciones", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, cauciones)
}

// Update handles PUT /cauciones/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	req, err := decodeJSON[UpdateCaucionRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.uc.Update(r.Context(), id, req.toInput())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "caucion not found")
			return
		}
		h.log.Error("handler: update caucion", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, c)
}

// Delete handles DELETE /cauciones/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "caucion not found")
			return
		}
		h.log.Error("handler: delete caucion", "id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ChangeState handles PATCH /cauciones/{id}/estado.
func (h *Handler) ChangeState(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	req, err := decodeJSON[ChangeStateRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.uc.ChangeState(r.Context(), id, domain.Estado(req.Estado))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "caucion not found")
		case errors.Is(err, domain.ErrInvalidTransition):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrInvalidState):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.log.Error("handler: change state caucion", "id", id, "error", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, c)
}
