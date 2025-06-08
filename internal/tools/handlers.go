package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/piercegov/llm-npc-backend/internal/api"
	"github.com/piercegov/llm-npc-backend/internal/llm"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

// ToolHandlers contains all tool-related HTTP handlers
type ToolHandlers struct {
	sessionManager *SessionManager
}

// NewToolHandlers creates a new instance of tool handlers
func NewToolHandlers(sessionManager *SessionManager) *ToolHandlers {
	return &ToolHandlers{
		sessionManager: sessionManager,
	}
}

// RegisterHandler handles POST /tools/register
func (h *ToolHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req ToolRegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON", api.ErrCodeInvalidJSON, nil, r.Context())
		return
	}

	// Validate session ID
	if req.SessionID == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Session ID is required", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Validate tools
	if len(req.Tools) == 0 {
		api.WriteErrorResponse(w, http.StatusBadRequest, "At least one tool must be provided", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Validate each tool
	toolNames := make([]string, 0, len(req.Tools))
	for i, tool := range req.Tools {
		if tool.Name == "" {
			api.WriteErrorResponse(w, http.StatusBadRequest, "Tool name is required", api.ErrCodeValidation, map[string]string{"tool_index": fmt.Sprintf("%d", i)}, r.Context())
			return
		}
		if tool.Description == "" {
			api.WriteErrorResponse(w, http.StatusBadRequest, "Tool description is required", api.ErrCodeValidation, map[string]string{"tool_index": fmt.Sprintf("%d", i), "tool_name": tool.Name}, r.Context())
			return
		}

		// Validate parameters
		for paramName, param := range tool.Parameters {
			if param.Description == "" {
				api.WriteErrorResponse(w, http.StatusBadRequest, "Parameter description is required", api.ErrCodeValidation, map[string]string{
					"tool_name":      tool.Name,
					"parameter_name": paramName,
				}, r.Context())
				return
			}
			// Ensure valid parameter type
			switch param.Type {
			case llm.TypeString, llm.TypeNumber, llm.TypeBoolean, llm.TypeObject, llm.TypeArray:
				// Valid type
			default:
				api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid parameter type", api.ErrCodeValidation, map[string]string{
					"tool_name":      tool.Name,
					"parameter_name": paramName,
					"parameter_type": string(param.Type),
				}, r.Context())
				return
			}
		}

		toolNames = append(toolNames, tool.Name)
	}

	// Register the session and tools
	if err := h.sessionManager.RegisterSession(req.SessionID, req.Tools); err != nil {
		api.LogRequestError(r.Context(), "Failed to register tools", err)
		api.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to register tools", api.ErrCodeInternalServer, nil, r.Context())
		return
	}

	logging.Info("Tools registered successfully",
		"session_id", req.SessionID,
		"tools_count", len(req.Tools),
		"tool_names", toolNames,
	)

	response := ToolRegistrationResponse{
		SessionID:    req.SessionID,
		ToolsCount:   len(req.Tools),
		RegisteredAt: time.Now().UTC().Format(time.RFC3339),
		Success:      true,
		Message:      "Tools registered successfully",
		ToolNames:    toolNames,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// SessionInfoHandler handles GET /tools/session/{id}
func (h *ToolHandlers) SessionInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL path
	sessionID := r.URL.Path[len("/tools/session/"):]
	if sessionID == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Session ID is required", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Get session tools
	tools, err := h.sessionManager.GetSessionTools(sessionID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusNotFound, "Session not found", api.ErrCodeNotFound, nil, r.Context())
		return
	}

	// Extract tool names
	toolNames := make([]string, 0, len(tools))
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name)
	}

	response := map[string]interface{}{
		"session_id":  sessionID,
		"tools_count": len(tools),
		"tool_names":  toolNames,
		"tools":       tools,
		"success":     true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}