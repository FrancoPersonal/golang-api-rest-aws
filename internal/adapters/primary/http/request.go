package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/FrancoPersonal/golang-api-rest-aws/internal/domain/ports/primary"
)

// CreateCaucionRequest is the JSON body for POST /cauciones.
type CreateCaucionRequest struct {
	Numero           string    `json:"numero"`
	Tipo             string    `json:"tipo"`
	Monto            float64   `json:"monto"`
	Moneda           string    `json:"moneda"`
	FechaEmision     time.Time `json:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Beneficiario     string    `json:"beneficiario"`
	Tomador          string    `json:"tomador"`
}

func (r CreateCaucionRequest) validate() error {
	switch {
	case r.Numero == "":
		return errors.New("numero is required")
	case r.Tipo == "":
		return errors.New("tipo is required")
	case r.Monto <= 0:
		return errors.New("monto must be greater than 0")
	case r.Moneda == "":
		return errors.New("moneda is required")
	case r.Beneficiario == "":
		return errors.New("beneficiario is required")
	case r.Tomador == "":
		return errors.New("tomador is required")
	}
	return nil
}

func (r CreateCaucionRequest) toInput() primary.CreateCaucionInput {
	return primary.CreateCaucionInput{
		Numero:           r.Numero,
		Tipo:             r.Tipo,
		Monto:            r.Monto,
		Moneda:           r.Moneda,
		FechaEmision:     r.FechaEmision,
		FechaVencimiento: r.FechaVencimiento,
		Beneficiario:     r.Beneficiario,
		Tomador:          r.Tomador,
	}
}

// UpdateCaucionRequest is the JSON body for PUT /cauciones/{id}.
type UpdateCaucionRequest struct {
	Numero           *string    `json:"numero,omitempty"`
	Tipo             *string    `json:"tipo,omitempty"`
	Monto            *float64   `json:"monto,omitempty"`
	Moneda           *string    `json:"moneda,omitempty"`
	FechaEmision     *time.Time `json:"fecha_emision,omitempty"`
	FechaVencimiento *time.Time `json:"fecha_vencimiento,omitempty"`
	Beneficiario     *string    `json:"beneficiario,omitempty"`
	Tomador          *string    `json:"tomador,omitempty"`
}

func (r UpdateCaucionRequest) toInput() primary.UpdateCaucionInput {
	return primary.UpdateCaucionInput{
		Numero:           r.Numero,
		Tipo:             r.Tipo,
		Monto:            r.Monto,
		Moneda:           r.Moneda,
		FechaEmision:     r.FechaEmision,
		FechaVencimiento: r.FechaVencimiento,
		Beneficiario:     r.Beneficiario,
		Tomador:          r.Tomador,
	}
}

// ChangeStateRequest is the JSON body for PATCH /cauciones/{id}/estado.
type ChangeStateRequest struct {
	Estado string `json:"estado"`
}

func decodeJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}
