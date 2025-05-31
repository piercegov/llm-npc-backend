package api

import (
	"testing"

	"github.com/piercegov/llm-npc-backend/internal/logging"
)

func init() {
	// Initialize the logger for tests
	logging.InitLogger("debug")
}

// TestMain is used to set up testing environment
func TestMain(m *testing.M) {
	// Run the tests
	m.Run()
}