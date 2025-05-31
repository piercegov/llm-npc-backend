package api

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

// PanicRecoveryMiddleware catches and handles panics in HTTP handlers
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				stackTrace := debug.Stack()
				logging.Error("Panic recovered in HTTP handler",
					"error", err,
					"request_id", GetRequestID(r.Context()),
					"path", r.URL.Path,
					"method", r.Method,
					"stack_trace", string(stackTrace),
				)

				// Return a generic 500 error to not expose internal details
				WriteErrorResponse(
					w,
					http.StatusInternalServerError,
					"The server encountered an unexpected error and was unable to complete your request",
					ErrCodeInternalServer,
					nil,
					r.Context(),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequestTracingMiddleware adds a unique request ID to each request and logs timing information
func RequestTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate unique request ID
		requestID := uuid.New().String()

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Record start time
		startTime := time.Now()

		// Log incoming request
		logging.Info("Request started",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		// Create a wrapper for the ResponseWriter to capture the status code
		ww := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 if WriteHeader is never called
		}

		// Process the request
		next.ServeHTTP(ww, r)

		// Calculate request duration
		duration := time.Since(startTime)

		// Log request completion with status code and duration
		if ww.statusCode >= 400 {
			logging.Warn("Request completed with error",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", ww.statusCode,
				"duration_ms", duration.Milliseconds(),
			)
		} else {
			logging.Info("Request completed",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", ww.statusCode,
				"duration_ms", duration.Milliseconds(),
			)
		}
	})
}

// ValidationMiddleware provides HTTP method and content-type validation
func ValidationMiddleware(allowedMethods []string, requireJSON bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate HTTP method
			methodAllowed := false
			for _, method := range allowedMethods {
				if r.Method == method {
					methodAllowed = true
					break
				}
			}

			if !methodAllowed {
				// Set Allow header with accepted methods
				w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
				
				WriteErrorResponse(
					w,
					http.StatusMethodNotAllowed,
					fmt.Sprintf("Method %s not allowed, supported methods: %s", 
						r.Method, strings.Join(allowedMethods, ", ")),
					ErrCodeMethodNotAllowed,
					nil,
					r.Context(),
				)
				return
			}

			// Validate Content-Type for POST/PUT/PATCH requests requiring JSON
			if requireJSON && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
				contentType := r.Header.Get("Content-Type")
				if !strings.Contains(contentType, "application/json") {
					WriteErrorResponse(
						w,
						http.StatusUnsupportedMediaType,
						"Content-Type must be application/json",
						ErrCodeUnsupportedMedia,
						nil,
						r.Context(),
					)
					return
				}

				// If there's a body, try to parse it as JSON to validate early
				if r.ContentLength > 0 {
					var jsonBody map[string]interface{}
					bodyBytes, err := io.ReadAll(r.Body)
					if err != nil {
						WriteErrorResponse(
							w,
							http.StatusBadRequest,
							"Error reading request body",
							ErrCodeBadRequest,
							nil,
							r.Context(),
						)
						return
					}

					// Replace the body so it can be read again by handlers
					r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

					// Validate JSON
					if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
						WriteErrorResponse(
							w,
							http.StatusBadRequest,
							fmt.Sprintf("Invalid JSON format: %s", err.Error()),
							ErrCodeInvalidJSON,
							nil,
							r.Context(),
						)
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ErrorHandlingMiddleware provides standardized error handling utilities
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add error handling utilities to the request context
		ctx := context.WithValue(r.Context(), contextKey("error_handler"), true)
		r = r.WithContext(ctx)
		
		next.ServeHTTP(w, r)
	})
}

// statusResponseWriter is a wrapper around http.ResponseWriter to capture status code
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

// WriteHeader captures the status code before writing it
func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the number of bytes written
func (w *statusResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += int64(n)
	return n, err
}

// Unwrap returns the original ResponseWriter
func (w *statusResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// Flush implements http.Flusher if the underlying ResponseWriter implements it
func (w *statusResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker if the underlying ResponseWriter implements it
func (w *statusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("underlying ResponseWriter does not implement http.Hijacker")
}

// Push implements http.Pusher if the underlying ResponseWriter implements it
func (w *statusResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return errors.New("underlying ResponseWriter does not implement http.Pusher")
}

// ChainMiddleware creates a middleware chain in the correct order
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// WithMethodValidation is a helper for common method validation
func WithMethodValidation(handler http.HandlerFunc, methods ...string) http.Handler {
	return ChainMiddleware(
		handler,
		ValidationMiddleware(methods, false),
	)
}

// WithJSONValidation is a helper for endpoints requiring JSON
func WithJSONValidation(handler http.HandlerFunc, methods ...string) http.Handler {
	return ChainMiddleware(
		handler,
		ValidationMiddleware(methods, true),
	)
}

// ApplyDefaultMiddleware applies all the standard middleware in the correct order
func ApplyDefaultMiddleware(handler http.Handler) http.Handler {
	return ChainMiddleware(
		handler,
		RequestTracingMiddleware,
		PanicRecoveryMiddleware,
		ErrorHandlingMiddleware,
	)
}