package tools

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/piercegov/llm-npc-backend/internal/llm"
)

// ScratchpadStorage manages the persistent memory storage for all NPCs
type ScratchpadStorage struct {
	storage map[string]*NPCScratchpad
	mu      sync.RWMutex
}

// NPCScratchpad stores memory entries for a specific NPC
type NPCScratchpad struct {
	Entries map[string]ScratchpadEntry
	mu      sync.RWMutex
}

// ScratchpadEntry represents a single memory entry
type ScratchpadEntry struct {
	Value     string
	Timestamp time.Time
}

// NewScratchpadStorage creates a new scratchpad storage
func NewScratchpadStorage() *ScratchpadStorage {
	return &ScratchpadStorage{
		storage: make(map[string]*NPCScratchpad),
	}
}

// RegisterScratchpadTools registers all scratchpad-related tools
func RegisterScratchpadTools(registry *ToolRegistry, storage *ScratchpadStorage) error {
	// Write tool
	writeToolDef := llm.Tool{
		Name:        "write_scratchpad",
		Description: "Store a memory with a key and value in your persistent scratchpad",
		Parameters: map[string]llm.ToolParameter{
			"key": {
				Type:        llm.TypeString,
				Description: "The key to store the memory under",
				Required:    true,
			},
			"value": {
				Type:        llm.TypeString,
				Description: "The value to store",
				Required:    true,
			},
		},
	}
	
	if err := registry.RegisterTool(writeToolDef, storage.handleWrite); err != nil {
		return err
	}
	
	// Read tool
	readToolDef := llm.Tool{
		Name:        "read_scratchpad",
		Description: "Retrieve a specific memory by its key from your scratchpad",
		Parameters: map[string]llm.ToolParameter{
			"key": {
				Type:        llm.TypeString,
				Description: "The key of the memory to retrieve",
				Required:    true,
			},
		},
	}
	
	if err := registry.RegisterTool(readToolDef, storage.handleRead); err != nil {
		return err
	}
	
	// List tool
	listToolDef := llm.Tool{
		Name:        "list_scratchpad",
		Description: "List all memories stored in your scratchpad",
		Parameters:  map[string]llm.ToolParameter{}, // No parameters
	}
	
	if err := registry.RegisterTool(listToolDef, storage.handleList); err != nil {
		return err
	}
	
	// Delete tool
	deleteToolDef := llm.Tool{
		Name:        "delete_scratchpad",
		Description: "Delete a specific memory by its key from your scratchpad",
		Parameters: map[string]llm.ToolParameter{
			"key": {
				Type:        llm.TypeString,
				Description: "The key of the memory to delete",
				Required:    true,
			},
		},
	}
	
	if err := registry.RegisterTool(deleteToolDef, storage.handleDelete); err != nil {
		return err
	}
	
	return nil
}

// GetAllScratchpads returns all scratchpads for admin/debug purposes
func (s *ScratchpadStorage) GetAllScratchpads() map[string]map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := make(map[string]map[string]interface{})
	
	for npcID, scratchpad := range s.storage {
		scratchpad.mu.RLock()
		
		npcData := make(map[string]interface{})
		entries := make([]map[string]interface{}, 0, len(scratchpad.Entries))
		
		for key, entry := range scratchpad.Entries {
			entries = append(entries, map[string]interface{}{
				"key":       key,
				"value":     entry.Value,
				"timestamp": entry.Timestamp.Format(time.RFC3339),
			})
		}
		
		npcData["entries"] = entries
		npcData["count"] = len(entries)
		result[npcID] = npcData
		
		scratchpad.mu.RUnlock()
	}
	
	return result
}

// getOrCreateScratchpad returns the scratchpad for an NPC, creating it if it doesn't exist
func (s *ScratchpadStorage) getOrCreateScratchpad(npcID string) *NPCScratchpad {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	scratchpad, exists := s.storage[npcID]
	if !exists {
		scratchpad = &NPCScratchpad{
			Entries: make(map[string]ScratchpadEntry),
		}
		s.storage[npcID] = scratchpad
	}
	return scratchpad
}

// handleWrite handles the write_scratchpad tool
func (s *ScratchpadStorage) handleWrite(ctx context.Context, npcID string, args map[string]interface{}) (ToolResult, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return ToolResult{Success: false, Message: "key must be a non-empty string"}, fmt.Errorf("invalid key")
	}
	
	value, ok := args["value"].(string)
	if !ok {
		return ToolResult{Success: false, Message: "value must be a string"}, fmt.Errorf("invalid value")
	}
	
	scratchpad := s.getOrCreateScratchpad(npcID)
	
	scratchpad.mu.Lock()
	defer scratchpad.mu.Unlock()
	
	scratchpad.Entries[key] = ScratchpadEntry{
		Value:     value,
		Timestamp: time.Now(),
	}
	
	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Stored memory: %s = %s", key, value),
		Data: map[string]interface{}{
			"key":   key,
			"value": value,
		},
	}, nil
}

// handleRead handles the read_scratchpad tool
func (s *ScratchpadStorage) handleRead(ctx context.Context, npcID string, args map[string]interface{}) (ToolResult, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return ToolResult{Success: false, Message: "key must be a non-empty string"}, fmt.Errorf("invalid key")
	}
	
	s.mu.RLock()
	scratchpad, exists := s.storage[npcID]
	s.mu.RUnlock()
	
	if !exists {
		return ToolResult{Success: false, Message: fmt.Sprintf("No memory found with key: %s", key)}, nil
	}
	
	scratchpad.mu.RLock()
	defer scratchpad.mu.RUnlock()
	
	entry, exists := scratchpad.Entries[key]
	if !exists {
		return ToolResult{Success: false, Message: fmt.Sprintf("No memory found with key: %s", key)}, nil
	}
	
	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("%s: %s", key, entry.Value),
		Data: map[string]interface{}{
			"key":       key,
			"value":     entry.Value,
			"timestamp": entry.Timestamp.Format(time.RFC3339),
		},
	}, nil
}

// handleList handles the list_scratchpad tool
func (s *ScratchpadStorage) handleList(ctx context.Context, npcID string, args map[string]interface{}) (ToolResult, error) {
	s.mu.RLock()
	scratchpad, exists := s.storage[npcID]
	s.mu.RUnlock()
	
	if !exists || len(scratchpad.Entries) == 0 {
		return ToolResult{Success: true, Message: "No memories stored"}, nil
	}
	
	scratchpad.mu.RLock()
	defer scratchpad.mu.RUnlock()
	
	memories := make([]map[string]interface{}, 0, len(scratchpad.Entries))
	message := "Stored memories:\n"
	
	for key, entry := range scratchpad.Entries {
		memories = append(memories, map[string]interface{}{
			"key":       key,
			"value":     entry.Value,
			"timestamp": entry.Timestamp.Format(time.RFC3339),
		})
		message += fmt.Sprintf("- %s: %s\n", key, entry.Value)
	}
	
	return ToolResult{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"memories": memories,
			"count":    len(memories),
		},
	}, nil
}

// handleDelete handles the delete_scratchpad tool
func (s *ScratchpadStorage) handleDelete(ctx context.Context, npcID string, args map[string]interface{}) (ToolResult, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return ToolResult{Success: false, Message: "key must be a non-empty string"}, fmt.Errorf("invalid key")
	}
	
	s.mu.RLock()
	scratchpad, exists := s.storage[npcID]
	s.mu.RUnlock()
	
	if !exists {
		return ToolResult{Success: false, Message: fmt.Sprintf("No memory found with key: %s", key)}, nil
	}
	
	scratchpad.mu.Lock()
	defer scratchpad.mu.Unlock()
	
	if _, exists := scratchpad.Entries[key]; !exists {
		return ToolResult{Success: false, Message: fmt.Sprintf("No memory found with key: %s", key)}, nil
	}
	
	delete(scratchpad.Entries, key)
	
	return ToolResult{
		Success: true,
		Message: fmt.Sprintf("Deleted memory with key: %s", key),
		Data: map[string]interface{}{
			"key": key,
		},
	}, nil
}