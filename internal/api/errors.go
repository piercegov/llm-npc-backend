package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/piercegov/llm-npc-backend/internal/logging"
)

// ErrorResponse is the standard error response format for all API endpoints
type ErrorResponse struct {
	Error     string            `json:"error"`
	Code      string            `json:"code"`
	RequestID string            `json:"request_id,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

// Error codes for different categories
const (
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeInvalidJSON        = "INVALID_JSON"
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeMethodNotAllowed   = "METHOD_NOT_ALLOWED"
	ErrCodeUnsupportedMedia   = "UNSUPPORTED_MEDIA_TYPE"
	ErrCodeRateLimit          = "RATE_LIMIT_EXCEEDED"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeBadRequest         = "BAD_REQUEST"
	
	// LLM-specific error codes
	ErrCodeLLMProviderUnavailable = "LLM_PROVIDER_UNAVAILABLE"
	ErrCodeLLMBadRequest          = "LLM_BAD_REQUEST"
	ErrCodeLLMRateLimited         = "LLM_RATE_LIMITED"
	ErrCodeLLMTimeout             = "LLM_TIMEOUT"
	ErrCodeLLMUnauthorized        = "LLM_UNAUTHORIZED"
	ErrCodeLLMModelNotFound       = "LLM_MODEL_NOT_FOUND"
)

// Map HTTP status codes to error codes
var statusToErrorCode = map[int]string{
	http.StatusBadRequest:           ErrCodeBadRequest,
	http.StatusNotFound:             ErrCodeNotFound,
	http.StatusMethodNotAllowed:     ErrCodeMethodNotAllowed,
	http.StatusUnsupportedMediaType: ErrCodeUnsupportedMedia,
	http.StatusTooManyRequests:      ErrCodeRateLimit,
	http.StatusInternalServerError:  ErrCodeInternalServer,
	http.StatusServiceUnavailable:   ErrCodeServiceUnavailable,
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// RequestIDKey is the key used to store and retrieve the request ID from context
const RequestIDKey contextKey = "request_id"

// GetRequestID extracts the request ID from the context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

// LogRequestError logs an error with the request context and additional fields
func LogRequestError(ctx context.Context, message string, err error, additionalFields ...any) {
	fields := []any{"error", err}

	// Add request ID if available
	if reqID := GetRequestID(ctx); reqID != "" {
		fields = append(fields, "request_id", reqID)
	}

	// Add any additional fields
	fields = append(fields, additionalFields...)

	logging.Error(message, fields...)
}

// WriteErrorResponse writes a standardized error response to the http.ResponseWriter
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, code string, details map[string]string, ctx context.Context) {
	// If no specific code is provided, use the default mapping
	if code == "" {
		code = statusToErrorCode[statusCode]
		if code == "" {
			code = ErrCodeInternalServer
		}
	}

	// Create error response
	errorResp := ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}

	// Add request ID if available
	if ctx != nil {
		if reqID := GetRequestID(ctx); reqID != "" {
			errorResp.RequestID = reqID
		}
	}

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Write JSON response
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(errorResp); err != nil {
		logging.Error("Failed to encode error response", "error", err)
		// If JSON encoding fails, write a simple text response
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
