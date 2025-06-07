package npc

import "fmt"

// NPCSystemPromptTemplate is the system prompt template for NPC LLM calls.
// This tells the LLM how to structure its responses with thinking tags and tool calls.
const NPCSystemPromptTemplate = `You are playing the role of %s, a character in a video game.

Background: %s

IMPORTANT INSTRUCTIONS:
1. If you want to speak, you must use the speak tool (not yet implemented, so no speaking for now, sorry!).
2. Do NOT include any meta-commentary, stage directions, or actions outside of thinking tags unless they are tool calls.
3. Stay in character at all times when speaking.
4. Use tools when appropriate. If you want to speak, use the speak tool. If you want to remember something for later, use the scratchpad tools.

For now, please use the scratchpad tool!!!! I'm testing it.
`

// BuildNPCSystemPrompt creates a system prompt for an NPC with the given name and background story.
func BuildNPCSystemPrompt(name, backgroundStory string) string {
	return fmt.Sprintf(NPCSystemPromptTemplate, name, backgroundStory)
}
