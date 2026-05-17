package dto

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// ErrorBody contains the error code and message returned in failed responses.
type ErrorBody struct {
	Code    string `json:"code"    example:"invalid_input"`
	Message string `json:"message" example:"field 'numero' is required"`
}

// StandardResponse is the envelope used for all API responses.
type StandardResponse struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

func Success(statusCode int, data any) events.APIGatewayProxyResponse {
	return JSON(statusCode, StandardResponse{
		Success: true,
		Data:    data,
	})
}

func Fail(statusCode int, code, message string) events.APIGatewayProxyResponse {
	return JSON(statusCode, StandardResponse{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func JSON(statusCode int, payload StandardResponse) events.APIGatewayProxyResponse {
	body, err := json.Marshal(payload)
	if err != nil {
		fallback := `{"success":false,"error":{"code":"internal_error","message":"internal error"}}`
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fallback,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}
