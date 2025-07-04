package llm

import (
	"errors"
	"testing"
)

func TestProviderError(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		model    string
		baseErr  error
		message  string
		wantMsg  string
	}{
		{
			name:     "with message",
			provider: "ollama",
			model:    "qwen3:1.7b",
			baseErr:  ErrProviderUnavailable,
			message:  "service unavailable",
			wantMsg:  "ollama provider error (model: qwen3:1.7b): service unavailable - llm provider unavailable",
		},
		{
			name:     "without message",
			provider: "lmstudio",
			model:    "model",
			baseErr:  ErrProviderUnavailable,
			message:  "",
			wantMsg:  "lmstudio provider error (model: model): llm provider unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewProviderError(tt.provider, tt.model, tt.baseErr, tt.message)
			if err.Error() != tt.wantMsg {
				t.Errorf("ProviderError.Error() = %v, want %v", err.Error(), tt.wantMsg)
			}
			if !errors.Is(err, tt.baseErr) {
				t.Errorf("ProviderError does not wrap baseErr correctly")
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		want    bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "ErrProviderUnavailable",
			err:  ErrProviderUnavailable,
			want: true,
		},
		{
			name: "ErrTimeout",
			err:  ErrTimeout,
			want: true,
		},
		{
			name: "ErrRateLimited",
			err:  ErrRateLimited,
			want: true,
		},
		{
			name: "ErrBadRequest",
			err:  ErrBadRequest,
			want: false,
		},
		{
			name: "ProviderError with rate limited",
			err:  NewProviderError("test", "model", ErrRateLimited, "rate limited"),
			want: true,
		},
		{
			name: "ProviderError with unavailable",
			err:  NewProviderError("test", "model", ErrProviderUnavailable, "unavailable"),
			want: true,
		},
		{
			name: "ProviderError with bad request",
			err:  NewProviderError("test", "model", ErrBadRequest, "bad request"),
			want: false,
		},
		{
			name: "ProviderError with model not found",
			err:  NewProviderError("test", "model", ErrModelNotFound, "model not found"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.want {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}