package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Códigos de erro da SDK.
const (
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeRateLimit    = "RATE_LIMIT_EXCEEDED"
	ErrCodeInternal     = "INTERNAL_SERVER_ERROR"
	ErrCodeUnknown      = "UNKNOWN"
)

// Error representa um erro retornado pela API ou pela SDK.
type Error struct {
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Details    map[string]any `json:"details,omitempty"`
	HTTPStatus int            `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("agentsdk [%d %s]: %s", e.HTTPStatus, e.Code, e.Message)
}

// Helper functions para validação de erros

// IsNotFound retorna true se o erro for do tipo 404.
func IsNotFound(err error) bool { return isCode(err, ErrCodeNotFound) }

// IsUnauthorized retorna true se o erro for do tipo 401.
func IsUnauthorized(err error) bool { return isCode(err, ErrCodeUnauthorized) }

// IsForbidden retorna true se o erro for do tipo 403.
func IsForbidden(err error) bool { return isCode(err, ErrCodeForbidden) }

// IsRateLimit retorna true se o erro for de rate limit (429).
func IsRateLimit(err error) bool { return isCode(err, ErrCodeRateLimit) }

func isCode(err error, code string) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

// apiErrorBody é o formato esperado do corpo de erro da API.
type apiErrorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// ParseAPIError processa a resposta de erro da API e retorna um tipo Error estruturado.
func ParseAPIError(statusCode int, body []byte) *Error {
	e := &Error{HTTPStatus: statusCode, Code: httpStatusToCode(statusCode)}

	var apiErr apiErrorBody
	if err := json.Unmarshal(body, &apiErr); err == nil {
		if apiErr.Code != "" {
			e.Code = apiErr.Code
		}
		e.Message = apiErr.Message
		e.Details = apiErr.Details
	} else {
		e.Message = string(body)
	}

	if e.Message == "" {
		e.Message = http.StatusText(statusCode)
	}
	return e
}

func httpStatusToCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return ErrCodeBadRequest
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusConflict:
		return ErrCodeConflict
	case http.StatusTooManyRequests:
		return ErrCodeRateLimit
	case http.StatusInternalServerError:
		return ErrCodeInternal
	default:
		return ErrCodeUnknown
	}
}
