package npc

import (
	"fmt"
	"log"

	"github.com/piercegov/llm-npc-backend/internal/kg"
	"github.com/piercegov/llm-npc-backend/internal/llm"
)

type NPCTickInput struct {
	Surroundings        []Surrounding
	KnowledgeGraph      kg.KnowledgeGraph
	NPCState            NPCState
	KnowledgeGraphDepth int
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
		log.Fatalf("Error parsing surroundings: %v", err)
	}
	knowledgeGraphString, err := ParseKnowledgeGraph(input)
	if err != nil {
		log.Fatalf("Error parsing knowledge graph: %v", err)
	}

	llmRequest := llm.LLMRequest{
		Prompt: surroundingsString + "\n" + knowledgeGraphString,
	}

	llmResponse, err := CallLLM(llmRequest)
	if err != nil {
		log.Fatalf("Error calling LLM: %v", err)
	}

	fmt.Println(llmResponse.Response)
}

func CallLLM(input llm.LLMRequest) (llm.LLMResponse, error) {
	return llm.LLMResponse{}, nil
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
		depth = 1
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

func ParseNPCState(input NPCTickInput) (NPCState, error) {
	return NPCState{}, nil
}
