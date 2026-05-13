package domain

import "time"

// Estado represents the lifecycle state of a Caucion.
type Estado string

const (
	EstadoPendiente Estado = "pendiente"
	EstadoVigente   Estado = "vigente"
	EstadoVencida   Estado = "vencida"
	EstadoCancelada Estado = "cancelada"
)

// validTransitions defines the allowed state transitions.
var validTransitions = map[Estado][]Estado{
	EstadoPendiente: {EstadoVigente, EstadoCancelada},
	EstadoVigente:   {EstadoVencida, EstadoCancelada},
}

// CanTransitionTo reports whether transitioning from the current state to next is valid.
func (e Estado) CanTransitionTo(next Estado) bool {
	allowed, ok := validTransitions[e]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

// Caucion is the core domain entity.
type Caucion struct {
	ID               string    `json:"id"                dynamodbav:"id"`
	Numero           string    `json:"numero"            dynamodbav:"numero"`
	Tipo             string    `json:"tipo"              dynamodbav:"tipo"`
	Monto            float64   `json:"monto"             dynamodbav:"monto"`
	Moneda           string    `json:"moneda"            dynamodbav:"moneda"`
	FechaEmision     time.Time `json:"fecha_emision"     dynamodbav:"fecha_emision"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" dynamodbav:"fecha_vencimiento"`
	Estado           Estado    `json:"estado"            dynamodbav:"estado"`
	Beneficiario     string    `json:"beneficiario"      dynamodbav:"beneficiario"`
	Tomador          string    `json:"tomador"           dynamodbav:"tomador"`
	CreatedAt        time.Time `json:"created_at"        dynamodbav:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"        dynamodbav:"updated_at"`
}
