package npc

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	npcSystemPromptTemplate string
	promptOnce              sync.Once
)

// loadNPCSystemPrompt loads the NPC system prompt template from file
func loadNPCSystemPrompt() {
	promptOnce.Do(func() {
		promptPath := filepath.Join("prompts", "npc_system.txt")
		content, err := os.ReadFile(promptPath)
		if err != nil {
			// Fallback to embedded prompt if file doesn't exist
			npcSystemPromptTemplate = `You are playing the role of %s, a character in a video game.

Background: %s

IMPORTANT INSTRUCTIONS:
1. If you want to speak, you must use the speak tool.
2. Do NOT include any meta-commentary, stage directions, or actions outside of thinking tags unless they are tool calls.
3. Stay in character at all times when speaking.
4. Use tools when appropriate. If you want to speak, use the speak tool. If you want to remember something for later, use the scratchpad tools.

There is no actual user, think of the user as the game itself. You are a character in a video game. You are interacting with the world around you, as well as other characters.
You don't always need to do something. If you don't have anything to do, you can just think.`
			return
		}
		npcSystemPromptTemplate = string(content)
	})
}

// BuildNPCSystemPrompt creates a system prompt for an NPC with the given name and background story.
func BuildNPCSystemPrompt(name, backgroundStory string) string {
	loadNPCSystemPrompt()
	return fmt.Sprintf(npcSystemPromptTemplate, name, backgroundStory)
}
