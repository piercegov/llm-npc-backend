package kg

// TODO: KnowledgeGraph is currently just in memory, but will be persisted in SpacetimeDB in the future.
type KnowledgeGraph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

type Edge struct {
	Source string                 `json:"source"`
	Target string                 `json:"target"`
	Data   map[string]interface{} `json:"data"`
}
