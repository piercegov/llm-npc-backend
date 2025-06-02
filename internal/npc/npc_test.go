package npc

import (
	"fmt"
	"testing"

	"github.com/piercegov/llm-npc-backend/internal/kg"
)

func TestParseSurroundings(t *testing.T) {
	surroundings := []Surrounding{
		{Name: "Surrounding 1", Description: "Description 1"},
		{Name: "Surrounding 2", Description: "Description 2"},
	}

	surroundingsString, err := ParseSurroundings(NPCTickInput{Surroundings: surroundings})
	if err != nil {
		t.Errorf("Error parsing surroundings: %v", err)
	}

	fmt.Println(surroundingsString)
	expected := "<surroundings>\n\t<surrounding>\n\t\t<surrounding_name>Surrounding 1</surrounding_name>\n\t\t<surrounding_description>Description 1</surrounding_description>\n\t</surrounding>\n\t<surrounding>\n\t\t<surrounding_name>Surrounding 2</surrounding_name>\n\t\t<surrounding_description>Description 2</surrounding_description>\n\t</surrounding>\n</surroundings>"
	if surroundingsString != expected {
		t.Errorf("Expected %s, got %s", expected, surroundingsString)
	}
}

func TestParseKnowledgeGraph(t *testing.T) {
	knowledgeGraph := kg.KnowledgeGraph{
		Nodes: []kg.Node{
			{ID: "Node 1", Data: map[string]interface{}{"name": "Node 1"}},
		},
	}

	knowledgeGraphString, err := ParseKnowledgeGraph(NPCTickInput{KnowledgeGraph: knowledgeGraph})
	if err != nil {
		t.Errorf("Error parsing knowledge graph: %v", err)
	}

	fmt.Println(knowledgeGraphString)
	expected := "<knowledge_graph>\n\t<nodes>\n\t\t<node>\n\t\t\t<node_id>Node 1</node_id>\n\t\t\t<node_data>map[name:Node 1]</node_data>\n\t\t</node>\n\t</nodes>\n\t<edges>\n\t</edges>\n</knowledge_graph>"
	if knowledgeGraphString != expected {
		t.Errorf("Expected %s, got %s", expected, knowledgeGraphString)
	}
}
