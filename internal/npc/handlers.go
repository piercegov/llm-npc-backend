package npc

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/piercegov/llm-npc-backend/internal/api"
	"github.com/piercegov/llm-npc-backend/internal/logging"
	"github.com/piercegov/llm-npc-backend/internal/tools"
)

// NPCHandlers contains all NPC-related HTTP handlers
type NPCHandlers struct {
	storage      *NPCStorage
	toolRegistry *tools.ToolRegistry
}

// NewNPCHandlers creates a new instance of NPC handlers
func NewNPCHandlers(storage *NPCStorage, toolRegistry *tools.ToolRegistry) *NPCHandlers {
	return &NPCHandlers{
		storage:      storage,
		toolRegistry: toolRegistry,
	}
}

// RegisterHandler handles POST /npc/register
func (h *NPCHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req NPCRegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON", api.ErrCodeInvalidJSON, nil, r.Context())
		return
	}

	// Validate required fields
	if req.Name == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Name is required", api.ErrCodeValidation, nil, r.Context())
		return
	}
	if req.BackgroundStory == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Background story is required", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Register the NPC
	npcID, err := h.storage.Register(req.Name, req.BackgroundStory)
	if err != nil {
		api.LogRequestError(r.Context(), "Failed to register NPC", err)
		api.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to register NPC", api.ErrCodeInternalServer, nil, r.Context())
		return
	}

	logging.Info("NPC registered successfully", "npc_id", npcID, "name", req.Name)

	response := NPCRegisterResponse{
		NPCID:   npcID,
		Success: true,
		Message: "NPC registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ActHandler handles POST /npc/act
func (h *NPCHandlers) ActHandler(w http.ResponseWriter, r *http.Request) {
	var req NPCActRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON", api.ErrCodeInvalidJSON, nil, r.Context())
		return
	}

	// Validate NPC ID
	if req.NPCID == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "NPC ID is required", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Get the NPC
	npc, err := h.storage.Get(req.NPCID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusNotFound, "NPC not found", api.ErrCodeNotFound, nil, r.Context())
		return
	}

	// Set the tool registry in the input
	req.NPCTickInput.ToolRegistry = h.toolRegistry

	// Execute the tick
	result := npc.ActForTick(req.NPCTickInput)

	response := NPCActResponse{
		NPCID:         req.NPCID,
		NPCTickResult: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListHandler handles GET /npc/list
func (h *NPCHandlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	npcs := h.storage.List()

	// Convert to response format
	npcInfos := make(map[string]NPCInfo)
	for id, npc := range npcs {
		npcInfos[id] = NPCInfo{
			Name:            npc.Name,
			BackgroundStory: npc.BackgroundStory,
		}
	}

	response := NPCListResponse{
		NPCs:    npcInfos,
		Success: true,
		Count:   len(npcInfos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetHandler handles GET /npc/{id}
func (h *NPCHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	// Extract NPC ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/npc/")
	npcID := strings.Split(path, "/")[0]

	if npcID == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "NPC ID is required", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Get the NPC
	npc, err := h.storage.Get(npcID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusNotFound, "NPC not found", api.ErrCodeNotFound, nil, r.Context())
		return
	}

	response := NPCGetResponse{
		NPCID: npcID,
		NPC: NPCInfo{
			Name:            npc.Name,
			BackgroundStory: npc.BackgroundStory,
		},
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteHandler handles DELETE /npc/{id}
func (h *NPCHandlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Extract NPC ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/npc/")
	npcID := strings.Split(path, "/")[0]

	if npcID == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, "NPC ID is required", api.ErrCodeValidation, nil, r.Context())
		return
	}

	// Delete the NPC
	err := h.storage.Delete(npcID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusNotFound, "NPC not found", api.ErrCodeNotFound, nil, r.Context())
		return
	}

	logging.Info("NPC deleted successfully", "npc_id", npcID)

	response := NPCDeleteResponse{
		NPCID:   npcID,
		Success: true,
		Message: "NPC deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
