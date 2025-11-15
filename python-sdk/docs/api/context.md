# Context API Reference

The context module provides builders for creating NPC surroundings, events, and knowledge graphs.

## Surroundings

::: llm_npc.context.Surroundings
    options:
      show_root_heading: true
      show_source: false

## Event

::: llm_npc.context.Event
    options:
      show_root_heading: true
      show_source: false

## KnowledgeGraph

::: llm_npc.context.KnowledgeGraph
    options:
      show_root_heading: true
      show_source: false

## ContextBuilder

::: llm_npc.context.ContextBuilder
    options:
      show_root_heading: true
      show_source: false

## Usage Examples

### Building Surroundings

```python
from llm_npc import Surroundings

surroundings = Surroundings()
surroundings.add("Forest", "A dark, dense forest")
surroundings.add("Path", "A winding dirt path")
surroundings.add("Cottage", "A small wooden cottage")

# Use with NPC
response = npc.act(surroundings)

# Or chain methods
surroundings = (Surroundings()
    .add("Forest", "A dark forest")
    .add("Path", "A dirt path")
    .add("Cottage", "A wooden cottage"))
```

### Creating Events

```python
from llm_npc import Event

events = [
    Event("arrival", "You just arrived at the forest"),
    Event("sound", "You hear wolves howling"),
    Event("discovery", "You spot a cottage ahead")
]

response = npc.act(surroundings, events)
```

### Building Knowledge Graphs

```python
from llm_npc import KnowledgeGraph

kg = KnowledgeGraph()

# Add nodes
kg.add_node("player", type="person", name="Hero")
kg.add_node("dragon", type="creature", name="Smaug")
kg.add_node("treasure", type="object", value="immense")

# Add relationships
kg.add_edge("dragon", "treasure", relationship="guards")
kg.add_edge("player", "dragon", relationship="hunting")

# Use with NPC
response = npc.act(surroundings, events, knowledge_graph=kg)

# Or chain methods
kg = (KnowledgeGraph()
    .add_node("player", type="person")
    .add_node("dragon", type="creature")
    .add_edge("player", "dragon", relationship="hunting"))
```

### Using Context Builder

```python
from llm_npc import ContextBuilder, KnowledgeGraph

# Build complete context
context = ContextBuilder()
context.add_surrounding("Castle", "A grand castle")
context.add_surrounding("King", "The king on his throne")
context.add_event("arrival", "You were summoned")

# Optional: add knowledge graph
kg = KnowledgeGraph()
kg.add_node("king", disposition="pleased")
context.set_knowledge_graph(kg)

# Use with NPC
response = npc.act(context)

# Or chain all methods
context = (ContextBuilder()
    .add_surrounding("Castle", "A grand castle")
    .add_event("arrival", "You were summoned")
    .set_knowledge_graph(kg))
```

### Accessing Surroundings

```python
surroundings = Surroundings()
surroundings.add("Forest", "Dark forest")
surroundings.add("Sword", "Rusty sword")

# Get count
print(len(surroundings))  # 2

# Iterate
for item in surroundings:
    print(f"{item.name}: {item.description}")

# Index access
first = surroundings[0]
print(first.name)  # "Forest"

# Convert to list
items_list = surroundings.to_list()
```

### Knowledge Graph Inspection

```python
kg = KnowledgeGraph()
kg.add_node("a", type="person")
kg.add_node("b", type="person")
kg.add_edge("a", "b", relationship="friend")

# Get size (nodes + edges)
print(len(kg))  # 3

# Convert to dict
graph_dict = kg.to_dict()
print(graph_dict["nodes"])
print(graph_dict["edges"])
```

## See Also

- [Building Context Guide](../user-guide/context.md)
- [Knowledge Graphs Guide](../user-guide/knowledge-graphs.md)
- [Models API](models.md)

