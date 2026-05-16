package domain

import "time"

const EstadoPendiente = "pendiente"

type Caucion struct {
	ID               string    `json:"id"`
	Numero           string    `json:"numero"`
	Tipo             string    `json:"tipo"`
	Monto            float64   `json:"monto"`
	Moneda           string    `json:"moneda"`
	Estado           string    `json:"estado"`
	Beneficiario     string    `json:"beneficiario"`
	Tomador          string    `json:"tomador"`
	FechaEmision     time.Time `json:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateCaucionInput struct {
	Numero           string
	Tipo             string
	Monto            float64
	Moneda           string
	Beneficiario     string
	Tomador          string
	FechaEmision     time.Time
	FechaVencimiento time.Time
}
