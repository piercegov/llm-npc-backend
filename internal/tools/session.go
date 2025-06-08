package tools

import (
	"fmt"
	"sync"
	"time"

	"github.com/piercegov/llm-npc-backend/internal/llm"
)

// Session represents a game session with its custom tools
type Session struct {
	ID        string
	Tools     map[string]llm.Tool
	CreatedAt time.Time
	LastUsed  time.Time
}

// SessionManager manages game sessions and their custom tools
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	// Configuration
	expirationDuration time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions:           make(map[string]*Session),
		expirationDuration: 1 * time.Hour, // Sessions expire after 1 hour of inactivity
	}

	// Start cleanup goroutine
	go sm.cleanupExpiredSessions()

	return sm
}

// RegisterSession creates or updates a session with custom tools
func (sm *SessionManager) RegisterSession(sessionID string, tools []llm.Tool) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		session = &Session{
			ID:        sessionID,
			Tools:     make(map[string]llm.Tool),
			CreatedAt: time.Now(),
			LastUsed:  time.Now(),
		}
		sm.sessions[sessionID] = session
	} else {
		session.LastUsed = time.Now()
	}

	// Add or update tools for this session
	for _, tool := range tools {
		session.Tools[tool.Name] = tool
	}

	return nil
}

// GetSessionTools returns all tools registered for a session
func (sm *SessionManager) GetSessionTools(sessionID string) ([]llm.Tool, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Update last used time
	session.LastUsed = time.Now()

	// Convert map to slice
	tools := make([]llm.Tool, 0, len(session.Tools))
	for _, tool := range session.Tools {
		tools = append(tools, tool)
	}

	return tools, nil
}

// TouchSession updates the last used time for a session
func (sm *SessionManager) TouchSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[sessionID]; exists {
		session.LastUsed = time.Now()
	}
}

// DeleteSession removes a session and its tools
func (sm *SessionManager) DeleteSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sessions[sessionID]; !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	delete(sm.sessions, sessionID)
	return nil
}

// GetSessionCount returns the number of active sessions
func (sm *SessionManager) GetSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// cleanupExpiredSessions runs periodically to remove expired sessions
func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for sessionID, session := range sm.sessions {
			if now.Sub(session.LastUsed) > sm.expirationDuration {
				delete(sm.sessions, sessionID)
			}
		}
		sm.mu.Unlock()
	}
}

// ToolRegistrationRequest represents the request to register tools for a session
type ToolRegistrationRequest struct {
	SessionID string     `json:"session_id" binding:"required"`
	Tools     []llm.Tool `json:"tools" binding:"required"`
}

// ToolRegistrationResponse represents the response from tool registration
type ToolRegistrationResponse struct {
	SessionID    string   `json:"session_id"`
	ToolsCount   int      `json:"tools_count"`
	RegisteredAt string   `json:"registered_at"`
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	ToolNames    []string `json:"tool_names"`
}