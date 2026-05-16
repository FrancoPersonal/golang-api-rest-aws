package dto

import "time"

type CreateSuretyBondRequest struct {
	Number      string    `json:"numero"`
	Type        string    `json:"tipo"`
	Amount      float64   `json:"monto"`
	Currency    string    `json:"moneda"`
	Beneficiary string    `json:"beneficiario"`
	Holder      string    `json:"tomador"`
	IssueDate   time.Time `json:"fecha_emision"`
	ExpiryDate  time.Time `json:"fecha_vencimiento"`
}
