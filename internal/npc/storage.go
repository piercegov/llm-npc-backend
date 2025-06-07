package npc

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// NPCStorage provides thread-safe in-memory storage for NPCs
type NPCStorage struct {
	npcs map[string]*NPC
	mu   sync.RWMutex
}

// NewNPCStorage creates a new NPC storage instance
func NewNPCStorage() *NPCStorage {
	return &NPCStorage{
		npcs: make(map[string]*NPC),
	}
}

// Register adds a new NPC and returns its generated ID
func (s *NPCStorage) Register(name, backgroundStory string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate unique ID
	id := uuid.New().String()

	// Create NPC
	npc := &NPC{
		Name:            name,
		BackgroundStory: backgroundStory,
	}

	// Store NPC
	s.npcs[id] = npc

	return id, nil
}

// Get retrieves an NPC by ID
func (s *NPCStorage) Get(id string) (*NPC, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	npc, exists := s.npcs[id]
	if !exists {
		return nil, fmt.Errorf("NPC with ID %s not found", id)
	}

	return npc, nil
}

// List returns all registered NPCs with their IDs
func (s *NPCStorage) List() map[string]*NPC {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]*NPC)
	for id, npc := range s.npcs {
		result[id] = npc
	}

	return result
}

// Delete removes an NPC by ID
func (s *NPCStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.npcs[id]; !exists {
		return fmt.Errorf("NPC with ID %s not found", id)
	}

	delete(s.npcs, id)
	return nil
}

// Count returns the number of registered NPCs
func (s *NPCStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.npcs)
}