package dto

import "time"

// CreateSuretyBondRequest contains the fields required to create a surety bond.
type CreateSuretyBondRequest struct {
	Number      string    `json:"numero"      example:"CB-2024-001"`
	Type        string    `json:"tipo"        example:"performance"`
	Amount      float64   `json:"monto"       example:"50000"`
	Currency    string    `json:"moneda"      example:"USD"`
	Beneficiary string    `json:"beneficiario" example:"ACME Corp"`
	Holder      string    `json:"tomador"     example:"John Doe"`
	IssueDate   time.Time `json:"fecha_emision" example:"2024-01-15T10:00:00Z"`
	ExpiryDate  time.Time `json:"fecha_vencimiento" example:"2025-01-15T10:00:00Z"`
}
