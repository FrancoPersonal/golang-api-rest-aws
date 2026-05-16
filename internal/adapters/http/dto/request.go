package dto

import "time"

type CreateCaucionRequest struct {
	Numero           string    `json:"numero"`
	Tipo             string    `json:"tipo"`
	Monto            float64   `json:"monto"`
	Moneda           string    `json:"moneda"`
	Beneficiario     string    `json:"beneficiario"`
	Tomador          string    `json:"tomador"`
	FechaEmision     time.Time `json:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
}
