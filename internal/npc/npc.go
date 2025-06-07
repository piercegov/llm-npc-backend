package npc

import (
	"context"
	"fmt"
	"strings"

	"github.com/piercegov/llm-npc-backend/internal/kg"
	"github.com/piercegov/llm-npc-backend/internal/llm"
	"github.com/piercegov/llm-npc-backend/internal/logging"
	"github.com/piercegov/llm-npc-backend/internal/tools"
)

// Request/Response types for API endpoints

// NPCRegisterRequest represents the request to register a new NPC
type NPCRegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	BackgroundStory string `json:"background_story" binding:"required"`
}

// NPCRegisterResponse represents the response from registering an NPC
type NPCRegisterResponse struct {
	NPCID   string `json:"npc_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NPCActRequest represents the request to make an NPC act
type NPCActRequest struct {
	NPCID string `json:"npc_id" binding:"required"`
	NPCTickInput
}

// NPCActResponse represents the response from an NPC action
type NPCActResponse struct {
	NPCID string `json:"npc_id"`
	NPCTickResult
}

// NPCListResponse represents the response from listing NPCs
type NPCListResponse struct {
	NPCs    map[string]NPCInfo `json:"npcs"`
	Success bool               `json:"success"`
	Count   int                `json:"count"`
}

// NPCInfo represents basic NPC information for listing
type NPCInfo struct {
	Name            string `json:"name"`
	BackgroundStory string `json:"background_story"`
}

// NPCGetResponse represents the response from getting a specific NPC
type NPCGetResponse struct {
	NPCID   string  `json:"npc_id"`
	NPC     NPCInfo `json:"npc"`
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
}

// NPCDeleteResponse represents the response from deleting an NPC
type NPCDeleteResponse struct {
	NPCID   string `json:"npc_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type NPCTickResult struct {
	Rounds       []InferenceRound
	LLMResponse  string // Concatenated responses from all rounds
	Success      bool
	ErrorMessage string
}

type InferenceRound struct {
	RoundNumber  int
	LLMResponse  string
	ToolsUsed    []ToolResult
	Success      bool
	ErrorMessage string
}

type ToolResult struct {
	ToolName string
	Args     map[string]interface{}
	Success  bool
	Response string
	Error    string
}

// Maybe this should be related to the knowledge graph?
type NPCTickEvent struct {
	EventType        string
	EventDescription string
}

type NPCTickInput struct {
	Surroundings        []Surrounding
	KnowledgeGraph      kg.KnowledgeGraph
	NPCState            NPCState
	KnowledgeGraphDepth int
	Events              []NPCTickEvent
	ToolRegistry        *tools.ToolRegistry // Optional: if nil, no tools available
}
type NPC struct {
	Name            string
	BackgroundStory string
}

type NPCState struct {
	// NOTE: This should probably be configurable, like tools.
	// e.g. some games need health, some need inventory, faction, etc.
}

type Surrounding struct {
	Name        string
	Description string
}

func (n *NPC) ActForTick(input NPCTickInput) NPCTickResult {
	return n.actForTickWithDepth(input, 0)
}

const maxThinkingDepth = 3

func (n *NPC) actForTickWithDepth(input NPCTickInput, depth int) NPCTickResult {
	surroundingsString, err := ParseSurroundings(input)
	if err != nil {
		logging.Error("Error parsing surroundings: %v", err)
		return NPCTickResult{Success: false, ErrorMessage: fmt.Sprintf("Error parsing surroundings: %v", err)}
	}
	knowledgeGraphString, err := ParseKnowledgeGraph(input)
	if err != nil {
		logging.Error("Error parsing knowledge graph: %v", err)
		return NPCTickResult{Success: false, ErrorMessage: fmt.Sprintf("Error parsing knowledge graph: %v", err)}
	}
	eventsString, err := ParseEvents(input)
	if err != nil {
		logging.Error("Error parsing events: %v", err)
		return NPCTickResult{Success: false, ErrorMessage: fmt.Sprintf("Error parsing events: %v", err)}
	}

	systemPrompt := BuildNPCSystemPrompt(n.Name, n.BackgroundStory)

	// Log NPC action details
	logging.Info("NPC ActForTick",
		"npc_name", n.Name,
		"depth", depth,
		"surroundings_count", len(input.Surroundings),
		"events_count", len(input.Events),
		"knowledge_graph_nodes", len(input.KnowledgeGraph.Nodes),
		"knowledge_graph_edges", len(input.KnowledgeGraph.Edges),
	)

	llmRequest := llm.LLMRequest{
		SystemPrompt: systemPrompt,
		Prompt:       surroundingsString + "\n" + knowledgeGraphString + "\n" + eventsString,
	}

	// Add tools if available
	if input.ToolRegistry != nil {
		llmRequest.Tools = input.ToolRegistry.GetTools()
	}

	llmResponse, err := CallLLM(llmRequest)
	if err != nil {
		logging.Error("Error calling LLM: %v", err)
		return NPCTickResult{Success: false, ErrorMessage: fmt.Sprintf("Error calling LLM: %v", err)}
	}

	var toolResults []ToolResult

	// Process any tool uses
	if len(llmResponse.ToolUses) > 0 && input.ToolRegistry != nil {
		ctx := context.Background()
		for _, toolUse := range llmResponse.ToolUses {
			logging.Info("NPC using tool",
				"npc_name", n.Name,
				"tool_name", toolUse.ToolName,
				"args", toolUse.ToolArgs,
			)

			result, err := input.ToolRegistry.ExecuteTool(ctx, n.Name, toolUse)

			toolResult := ToolResult{
				ToolName: toolUse.ToolName,
				Args:     toolUse.ToolArgs,
				Success:  result.Success,
				Response: result.Message,
			}

			if err != nil {
				logging.Error("Error executing tool: %v", err)
				toolResult.Success = false
				toolResult.Error = err.Error()
			} else {
				logging.Info("Tool execution completed",
					"tool_name", toolUse.ToolName,
					"success", result.Success,
					"message", result.Message,
				)
			}

			toolResults = append(toolResults, toolResult)
		}
	}

	// Check if continue_thinking was used and we haven't exceeded depth limit
	var usedContinueThinking bool
	for _, toolUse := range llmResponse.ToolUses {
		if toolUse.ToolName == "continue_thinking" {
			usedContinueThinking = true
			break
		}
	}

	if usedContinueThinking && depth < maxThinkingDepth && input.ToolRegistry != nil {
		logging.Info("NPC continuing thinking",
			"npc_name", n.Name,
			"current_depth", depth,
		)

		// Convert tool results to events for the next thinking round
		var newEvents []NPCTickEvent
		for _, toolResult := range toolResults {
			eventType := "tool_execution"
			if !toolResult.Success {
				eventType = "tool_error"
			}

			description := fmt.Sprintf("Tool '%s' executed", toolResult.ToolName)
			if toolResult.Response != "" {
				description += fmt.Sprintf(" - Response: %s", toolResult.Response)
			}
			if toolResult.Error != "" {
				description += fmt.Sprintf(" - Error: %s", toolResult.Error)
			}

			newEvents = append(newEvents, NPCTickEvent{
				EventType:        eventType,
				EventDescription: description,
			})
		}

		// Create new input with tool results as events
		continueInput := NPCTickInput{
			Surroundings:        input.Surroundings,
			KnowledgeGraph:      input.KnowledgeGraph,
			NPCState:            input.NPCState,
			KnowledgeGraphDepth: input.KnowledgeGraphDepth,
			Events:              newEvents,
			ToolRegistry:        input.ToolRegistry,
		}

		// Recursively call for continued thinking
		continueResult := n.actForTickWithDepth(continueInput, depth+1)

		// Create current round
		currentRound := InferenceRound{
			RoundNumber: depth + 1,
			LLMResponse: llmResponse.Response,
			ToolsUsed:   toolResults,
			Success:     true,
		}

		// Combine rounds from current and recursive calls
		allRounds := append([]InferenceRound{currentRound}, continueResult.Rounds...)

		// Build concatenated response with inference markers
		var concatenatedResponse strings.Builder
		concatenatedResponse.WriteString(fmt.Sprintf("=== Inference %d ===\n%s", currentRound.RoundNumber, currentRound.LLMResponse))
		if continueResult.LLMResponse != "" {
			concatenatedResponse.WriteString("\n")
			concatenatedResponse.WriteString(continueResult.LLMResponse)
		}

		return NPCTickResult{
			Rounds:       allRounds,
			LLMResponse:  concatenatedResponse.String(),
			Success:      continueResult.Success,
			ErrorMessage: continueResult.ErrorMessage,
		}
	}

	// Base case - no continue_thinking used
	currentRound := InferenceRound{
		RoundNumber: depth + 1,
		LLMResponse: llmResponse.Response,
		ToolsUsed:   toolResults,
		Success:     true,
	}

	response := llmResponse.Response
	if depth > 0 {
		response = fmt.Sprintf("=== Inference %d ===\n%s", currentRound.RoundNumber, llmResponse.Response)
	}

	return NPCTickResult{
		Rounds:      []InferenceRound{currentRound},
		LLMResponse: response,
		Success:     true,
	}
}

func CallLLM(input llm.LLMRequest) (llm.LLMResponse, error) {
	// TODO: This should be configurable to support multiple LLM providers
	ollama := llm.NewOllama("11434")
	return ollama.Generate(input)
}

func ParseSurroundings(input NPCTickInput) (string, error) {
	surroundingsString := "<surroundings>\n"
	for _, surrounding := range input.Surroundings {
		surroundingsString += fmt.Sprintf("\t<surrounding>\n\t\t<surrounding_name>%s</surrounding_name>\n\t\t<surrounding_description>%s</surrounding_description>\n\t</surrounding>\n", surrounding.Name, surrounding.Description)
	}
	surroundingsString += "</surroundings>"
	return surroundingsString, nil
}

func ParseKnowledgeGraph(input NPCTickInput) (string, error) {
	depth := input.KnowledgeGraphDepth
	if depth == 0 {
		return "<knowledge_graph></knowledge_graph>", nil
	}

	kgString := "<knowledge_graph>\n"
	kgString += fmt.Sprintf("\t<nodes>\n")
	for _, node := range input.KnowledgeGraph.Nodes {
		kgString += fmt.Sprintf("\t\t<node>\n\t\t\t<node_id>%s</node_id>\n\t\t\t<node_data>%s</node_data>\n\t\t</node>\n", node.ID, node.Data)
	}
	kgString += "\t</nodes>\n"
	kgString += fmt.Sprintf("\t<edges>\n")
	for _, edge := range input.KnowledgeGraph.Edges {
		kgString += fmt.Sprintf("\t\t<edge>\n\t\t\t<edge_source>%s</edge_source>\n\t\t\t<edge_target>%s</edge_target>\n\t\t\t<edge_data>%s</edge_data>\n\t\t</edge>\n", edge.Source, edge.Target, edge.Data)
	}
	kgString += fmt.Sprintf("\t</edges>\n")
	kgString += "</knowledge_graph>"

	return kgString, nil
}

func ParseEvents(input NPCTickInput) (string, error) {
	if len(input.Events) == 0 {
		return "<events_since_last_tick></events_since_last_tick>", nil
	}

	eventsString := "<events_since_last_tick>\n"
	for _, event := range input.Events {
		eventsString += fmt.Sprintf("\t<event>\n\t\t<event_type>%s</event_type>\n\t\t<event_description>%s</event_description>\n\t</event>\n", event.EventType, event.EventDescription)
	}
	eventsString += "</events_since_last_tick>"
	return eventsString, nil
}

func ParseNPCState(input NPCTickInput) (NPCState, error) {
	return NPCState{}, nil
}
