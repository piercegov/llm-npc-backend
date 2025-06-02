package kg

type KnowledgeGraph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	ID   string
	Data map[string]interface{}
}

type Edge struct {
	Source string
	Target string
	Data   map[string]interface{}
}
