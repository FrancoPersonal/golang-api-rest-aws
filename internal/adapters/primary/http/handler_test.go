package httphandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httphandler "github.com/FrancoPersonal/golang-api-rest-aws/internal/adapters/primary/http"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain"
	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/primary"
	pkglogger "github.com/FrancoPersonal/golang-api-rest-aws/pkg/logger"
)

// ── mocks ────────────────────────────────────────────────────────────────────

type mockUseCase struct{ mock.Mock }

func (m *mockUseCase) Create(ctx context.Context, input primary.CreateCaucionInput) (*domain.Caucion, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Caucion), args.Error(1)
}
func (m *mockUseCase) GetByID(ctx context.Context, id string) (*domain.Caucion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Caucion), args.Error(1)
}
func (m *mockUseCase) List(ctx context.Context) ([]domain.Caucion, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Caucion), args.Error(1)
}
func (m *mockUseCase) Update(ctx context.Context, id string, input primary.UpdateCaucionInput) (*domain.Caucion, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Caucion), args.Error(1)
}
func (m *mockUseCase) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockUseCase) ChangeState(ctx context.Context, id string, newState domain.Estado) (*domain.Caucion, error) {
	args := m.Called(ctx, id, newState)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Caucion), args.Error(1)
}

type mockLogger struct{}

func (l *mockLogger) Info(msg string, args ...any)      {}
func (l *mockLogger) Error(msg string, args ...any)     {}
func (l *mockLogger) Debug(msg string, args ...any)     {}
func (l *mockLogger) Warn(msg string, args ...any)      {}
func (l *mockLogger) With(args ...any) pkglogger.Logger { return l }

func newRouter(uc *mockUseCase) http.Handler {
	return httphandler.NewRouter(uc, &mockLogger{})
}

func sampleCaucion() *domain.Caucion {
	return &domain.Caucion{
		ID:               "test-id",
		Numero:           "CAU-001",
		Tipo:             "garantia",
		Monto:            100_000.0,
		Moneda:           "ARS",
		FechaEmision:     time.Now().UTC(),
		FechaVencimiento: time.Now().UTC().AddDate(1, 0, 0),
		Estado:           domain.EstadoPendiente,
		Beneficiario:     "Empresa SA",
		Tomador:          "Cliente SRL",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
}

func mustJSON(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// ── POST /cauciones ───────────────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("Create", mock.Anything, mock.AnythingOfType("primary.CreateCaucionInput")).
		Return(sampleCaucion(), nil)

	body := mustJSON(map[string]any{
		"numero": "CAU-001", "tipo": "garantia", "monto": 100000.0,
		"moneda": "ARS", "beneficiario": "Empresa SA", "tomador": "Cliente SRL",
		"fecha_emision": time.Now(), "fecha_vencimiento": time.Now().AddDate(1, 0, 0),
	})
	req := httptest.NewRequest(http.MethodPost, "/cauciones", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)
}

func TestCreate_BadBody(t *testing.T) {
	uc := &mockUseCase{}
	req := httptest.NewRequest(http.MethodPost, "/cauciones", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_ValidationError(t *testing.T) {
	uc := &mockUseCase{}
	body := mustJSON(map[string]any{"monto": 0}) // missing required fields
	req := httptest.NewRequest(http.MethodPost, "/cauciones", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_UseCaseError(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("Create", mock.Anything, mock.AnythingOfType("primary.CreateCaucionInput")).
		Return(nil, fmt.Errorf("unexpected error"))

	body := mustJSON(map[string]any{
		"numero": "CAU-001", "tipo": "garantia", "monto": 1.0,
		"moneda": "ARS", "beneficiario": "X", "tomador": "Y",
	})
	req := httptest.NewRequest(http.MethodPost, "/cauciones", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── GET /cauciones ────────────────────────────────────────────────────────────

func TestList_Success(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("List", mock.Anything).Return([]domain.Caucion{*sampleCaucion()}, nil)

	req := httptest.NewRequest(http.MethodGet, "/cauciones", nil)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestList_Error(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("List", mock.Anything).Return(nil, fmt.Errorf("db error"))

	req := httptest.NewRequest(http.MethodGet, "/cauciones", nil)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── GET /cauciones/{id} ───────────────────────────────────────────────────────

func TestGetByID_Success(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("GetByID", mock.Anything, "test-id").Return(sampleCaucion(), nil)

	req := httptest.NewRequest(http.MethodGet, "/cauciones/test-id", nil)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetByID_NotFound(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("GetByID", mock.Anything, "missing").Return(nil, domain.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/cauciones/missing", nil)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── PUT /cauciones/{id} ───────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("Update", mock.Anything, "test-id", mock.AnythingOfType("primary.UpdateCaucionInput")).
		Return(sampleCaucion(), nil)

	body := mustJSON(map[string]any{"numero": "CAU-999"})
	req := httptest.NewRequest(http.MethodPut, "/cauciones/test-id", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate_NotFound(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("Update", mock.Anything, "missing", mock.AnythingOfType("primary.UpdateCaucionInput")).
		Return(nil, domain.ErrNotFound)

	body := mustJSON(map[string]any{"numero": "X"})
	req := httptest.NewRequest(http.MethodPut, "/cauciones/missing", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── DELETE /cauciones/{id} ────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("Delete", mock.Anything, "test-id").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/cauciones/test-id", nil)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDelete_NotFound(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("Delete", mock.Anything, "missing").Return(domain.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/cauciones/missing", nil)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── PATCH /cauciones/{id}/estado ──────────────────────────────────────────────

func TestChangeState_Success(t *testing.T) {
	uc := &mockUseCase{}
	updated := sampleCaucion()
	updated.Estado = domain.EstadoVigente
	uc.On("ChangeState", mock.Anything, "test-id", domain.EstadoVigente).Return(updated, nil)

	body := mustJSON(map[string]string{"estado": "vigente"})
	req := httptest.NewRequest(http.MethodPatch, "/cauciones/test-id/estado", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChangeState_InvalidTransition(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("ChangeState", mock.Anything, "test-id", domain.EstadoVencida).
		Return(nil, domain.ErrInvalidTransition)

	body := mustJSON(map[string]string{"estado": "vencida"})
	req := httptest.NewRequest(http.MethodPatch, "/cauciones/test-id/estado", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestChangeState_InvalidState(t *testing.T) {
	uc := &mockUseCase{}
	uc.On("ChangeState", mock.Anything, "test-id", domain.Estado("unknown")).
		Return(nil, domain.ErrInvalidState)

	body := mustJSON(map[string]string{"estado": "unknown"})
	req := httptest.NewRequest(http.MethodPatch, "/cauciones/test-id/estado", body)
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangeState_BadBody(t *testing.T) {
	uc := &mockUseCase{}
	req := httptest.NewRequest(http.MethodPatch, "/cauciones/test-id/estado", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()

	newRouter(uc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
