package http

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIResponse is the standard success envelope.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// APIError is the standard error envelope.
type APIError struct {
	Success bool           `json:"success"`
	Error   ErrorBody      `json:"error"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func JSONSuccess(c *gin.Context, status int, message string, data any) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONError(c *gin.Context, status int, code, message string, details map[string]string) {
	c.JSON(status, APIError{
		Success: false,
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationDetails extracts field-level errors from go-playground/validator.
func ValidationDetails(err error) map[string]string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return map[string]string{"_": err.Error()}
	}

	details := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		details[fe.Field()] = validationMessage(fe)
	}
	return details
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return "must be >= " + fe.Param()
	case "max":
		return "must be <= " + fe.Param()
	case "uuid":
		return "must be a valid UUID"
	default:
		return "failed on '" + fe.Tag() + "' validation"
	}
}
