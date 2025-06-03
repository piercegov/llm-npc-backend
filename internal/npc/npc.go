package npc

import (
	"fmt"

	"github.com/piercegov/llm-npc-backend/internal/kg"
	"github.com/piercegov/llm-npc-backend/internal/llm"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

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

func (n *NPC) ActForTick(input NPCTickInput) {
	surroundingsString, err := ParseSurroundings(input)
	if err != nil {
		logging.Error("Error parsing surroundings: %v", err)
	}
	knowledgeGraphString, err := ParseKnowledgeGraph(input)
	if err != nil {
		logging.Error("Error parsing knowledge graph: %v", err)
	}
	eventsString, err := ParseEvents(input)
	if err != nil {
		logging.Error("Error parsing events: %v", err)
	}

	systemPrompt := BuildNPCSystemPrompt(n.Name, n.BackgroundStory)

	// Log NPC action details
	logging.Info("NPC ActForTick",
		"npc_name", n.Name,
		"surroundings_count", len(input.Surroundings),
		"events_count", len(input.Events),
		"knowledge_graph_nodes", len(input.KnowledgeGraph.Nodes),
		"knowledge_graph_edges", len(input.KnowledgeGraph.Edges),
	)

	llmRequest := llm.LLMRequest{
		SystemPrompt: systemPrompt,
		Prompt:       surroundingsString + "\n" + knowledgeGraphString + "\n" + eventsString,
	}

	llmResponse, err := CallLLM(llmRequest)
	if err != nil {
		logging.Error("Error calling LLM: %v", err)
	}

	fmt.Println(llmResponse.Response)
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
