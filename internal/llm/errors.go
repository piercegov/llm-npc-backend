package llm

import (
	"errors"
	"fmt"
)

// ErrProviderUnavailable indicates the LLM provider is unreachable
var ErrProviderUnavailable = errors.New("llm provider unavailable")

// ErrBadRequest indicates invalid request parameters or prompt
var ErrBadRequest = errors.New("bad request to llm provider")

// ErrRateLimited indicates the provider returned a rate limit error
var ErrRateLimited = errors.New("llm provider rate limited")

// ErrTimeout indicates the request timed out
var ErrTimeout = errors.New("llm request timeout")

// ErrUnauthorized indicates authentication/authorization failure
var ErrUnauthorized = errors.New("llm provider unauthorized")

// ErrModelNotFound indicates the requested model is not available
var ErrModelNotFound = errors.New("llm model not found")

// ProviderError wraps an error with additional context
type ProviderError struct {
	Provider string
	Model    string
	Err      error
	Message  string
}

func (e *ProviderError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s provider error (model: %s): %s - %v", e.Provider, e.Model, e.Message, e.Err)
	}
	return fmt.Sprintf("%s provider error (model: %s): %v", e.Provider, e.Model, e.Err)
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new provider error with context
func NewProviderError(provider, model string, err error, message string) *ProviderError {
	return &ProviderError{
		Provider: provider,
		Model:    model,
		Err:      err,
		Message:  message,
	}
}

// IsRetryable determines if an error is temporary and can be retried
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a wrapped error
	var provErr *ProviderError
	if errors.As(err, &provErr) {
		err = provErr.Err
	}
	
	// Rate limiting and timeouts are typically retryable
	return errors.Is(err, ErrRateLimited) || 
		errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrProviderUnavailable)
}