package domain

import "time"

const StatusPending = "pending"

type SuretyBond struct {
	ID          string    `json:"id"`
	Number      string    `json:"numero"`
	Type        string    `json:"tipo"`
	Amount      float64   `json:"monto"`
	Currency    string    `json:"moneda"`
	Status      string    `json:"estado"`
	Beneficiary string    `json:"beneficiario"`
	Holder      string    `json:"tomador"`
	IssueDate   time.Time `json:"fecha_emision"`
	ExpiryDate  time.Time `json:"fecha_vencimiento"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSuretyBondInput struct {
	Number      string
	Type        string
	Amount      float64
	Currency    string
	Beneficiary string
	Holder      string
	IssueDate   time.Time
	ExpiryDate  time.Time
}
