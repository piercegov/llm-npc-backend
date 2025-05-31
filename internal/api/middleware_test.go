package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPanicRecoveryMiddleware tests that panic recovery middleware catches panics
func TestPanicRecoveryMiddleware(t *testing.T) {
	// Create a handler that will panic
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap with panic recovery middleware
	handler := PanicRecoveryMiddleware(panicHandler)

	// Create a test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rec, req)

	// Check response
	res := rec.Result()
	defer res.Body.Close()

	// Should return 500 status code
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, res.StatusCode)
	}

	// Should return JSON with error response
	var errorResp ErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	// Check error code
	if errorResp.Code != ErrCodeInternalServer {
		t.Errorf("Expected error code %s, got %s", ErrCodeInternalServer, errorResp.Code)
	}
}

// TestRequestTracingMiddleware tests request ID generation and logging
func TestRequestTracingMiddleware(t *testing.T) {
	// Create a handler that extracts and returns the request ID
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		fmt.Fprint(w, requestID)
	})

	// Wrap with request tracing middleware
	tracedHandler := RequestTracingMiddleware(handler)

	// Create a test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Call the handler
	tracedHandler.ServeHTTP(rec, req)

	// Check response
	res := rec.Result()
	defer res.Body.Close()

	// Read the response body (which should contain the request ID)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %v", err)
	}
	requestID := string(body)

	// Check that we received a request ID
	if requestID == "" {
		t.Error("Expected a request ID, got empty string")
	}

	// Check that request ID header was set
	headerRequestID := res.Header.Get("X-Request-ID")
	if headerRequestID == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	// Header should match the ID returned in body
	if headerRequestID != requestID {
		t.Errorf("Request ID mismatch: header has %s, body has %s", headerRequestID, requestID)
	}
}

// TestValidationMiddleware_Methods tests HTTP method validation
func TestValidationMiddleware_Methods(t *testing.T) {
	// Create a simple handler
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with validation middleware that only allows GET
	handler := ValidationMiddleware([]string{"GET"}, false)(okHandler)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{"allowed method", "GET", http.StatusOK},
		{"disallowed method", "POST", http.StatusMethodNotAllowed},
		{"disallowed method", "PUT", http.StatusMethodNotAllowed},
		{"disallowed method", "DELETE", http.StatusMethodNotAllowed},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rec.Code)
			}

			if tc.expectedStatus == http.StatusMethodNotAllowed {
				// Check that Allow header is set correctly
				allowHeader := rec.Header().Get("Allow")
				if allowHeader != "GET" {
					t.Errorf("Expected Allow header to be 'GET', got '%s'", allowHeader)
				}

				// Should be JSON error response
				var errorResp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Could not decode response body: %v", err)
				}

				// Check error code
				if errorResp.Code != ErrCodeMethodNotAllowed {
					t.Errorf("Expected error code %s, got %s", ErrCodeMethodNotAllowed, errorResp.Code)
				}
			}
		})
	}
}

// TestValidationMiddleware_ContentType tests content type validation
func TestValidationMiddleware_ContentType(t *testing.T) {
	// Create a simple handler
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with validation middleware that requires JSON for POST
	handler := ValidationMiddleware([]string{"POST"}, true)(okHandler)

	tests := []struct {
		name           string
		contentType    string
		body           string
		expectedStatus int
	}{
		{"valid JSON", "application/json", `{"key":"value"}`, http.StatusOK},
		{"missing content type", "", `{"key":"value"}`, http.StatusUnsupportedMediaType},
		{"wrong content type", "text/plain", `{"key":"value"}`, http.StatusUnsupportedMediaType},
		{"invalid JSON", "application/json", `{"key":invalid}`, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rec.Code)
			}

			if tc.expectedStatus != http.StatusOK {
				// Should be JSON error response
				var errorResp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Could not decode response body: %v", err)
				}
			}
		})
	}
}

// TestChainMiddleware tests middleware chaining functionality
func TestChainMiddleware(t *testing.T) {
	executionOrder := []string{}

	// Create middleware functions that record execution order
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	// Create final handler
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Chain middleware: middleware1 -> middleware2 -> finalHandler
	chain := ChainMiddleware(finalHandler, middleware1, middleware2)

	// Create request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Execute the chain
	chain.ServeHTTP(rec, req)

	// Check execution order
	expected := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expected) {
		t.Fatalf("Expected %d middleware executions, got %d", len(expected), len(executionOrder))
	}

	for i, v := range expected {
		if executionOrder[i] != v {
			t.Errorf("Expected execution order at position %d to be '%s', got '%s'", 
				i, v, executionOrder[i])
		}
	}
}

// TestApplyDefaultMiddleware tests applying all default middleware
func TestApplyDefaultMiddleware(t *testing.T) {
	// Create a handler that panics to test recovery
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is in context
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Expected request ID in context")
		}
		
		// Trigger a panic to test recovery
		panic("test panic in default middleware chain")
	})

	// Apply default middleware
	handler := ApplyDefaultMiddleware(panicHandler)

	// Create request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(rec, req)

	// Should recover from panic and return 500
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d after panic recovery, got %d", 
			http.StatusInternalServerError, rec.Code)
	}

	// Should have request ID header
	if rec.Header().Get("X-Request-ID") == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	// Should return proper error response
	var errorResp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}
	
	if errorResp.Code != ErrCodeInternalServer {
		t.Errorf("Expected error code %s, got %s", 
			ErrCodeInternalServer, errorResp.Code)
	}
	
	// Request ID in the error response should match header
	if errorResp.RequestID != rec.Header().Get("X-Request-ID") {
		t.Errorf("Request ID mismatch: header has %s, response has %s", 
			rec.Header().Get("X-Request-ID"), errorResp.RequestID)
	}
}

// TestGetRequestID tests extracting request ID from context
func TestGetRequestID(t *testing.T) {
	// Test with nil context
	if id := GetRequestID(nil); id != "" {
		t.Errorf("Expected empty string for nil context, got '%s'", id)
	}
	
	// Test with context but no request ID
	if id := GetRequestID(httptest.NewRequest("GET", "/", nil).Context()); id != "" {
		t.Errorf("Expected empty string for context without request ID, got '%s'", id)
	}
}