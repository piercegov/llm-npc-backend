# Knowledge Graphs

Knowledge graphs give NPCs persistent memory and understanding of relationships between entities. They help NPCs remember past events, understand social connections, and make context-aware decisions.

## What is a Knowledge Graph?

A knowledge graph consists of:

- **Nodes**: Entities (people, places, objects, events, concepts)
- **Edges**: Relationships between nodes

Think of it as the NPC's memory and understanding of the world.

## Creating a Knowledge Graph

```python
from llm_npc import KnowledgeGraph

kg = KnowledgeGraph()

# Add nodes (entities)
kg.add_node("player", type="person", name="Hero", reputation="good")
kg.add_node("king", type="person", name="King Aldric", disposition="neutral")
kg.add_node("quest_01", type="quest", name="Slay the dragon", status="active")

# Add edges (relationships)
kg.add_edge("player", "king", relationship="serves")
kg.add_edge("player", "quest_01", relationship="accepted")
kg.add_edge("king", "quest_01", relationship="assigned")
```

## Using with NPCs

Pass the knowledge graph when NPCs act:

```python
response = npc.act(
    surroundings=["Throne room", "King Aldric on throne"],
    events=["You returned from the quest"],
    knowledge_graph=kg
)
```

The NPC now has context about the relationships and can reason about them.

## Node Structure

Nodes have an ID and arbitrary data:

```python
kg.add_node(
    "npc_merchant_01",           # Unique ID
    type="person",                # Data fields (arbitrary)
    name="Tobias the Merchant",
    occupation="merchant",
    location="Market Square",
    trust_level="medium",
    owes_money=True
)
```

### Common Node Types

```python
# People
kg.add_node("char_001", type="person", name="Alice", role="guard")

# Places
kg.add_node("loc_tavern", type="location", name="The Gilded Swan", area="downtown")

# Objects
kg.add_node("item_sword", type="object", name="Legendary Sword", rarity="epic")

# Events
kg.add_node("evt_battle", type="event", name="Battle of  Crossroads", date="3 days ago")

# Quests
kg.add_node("quest_dragon", type="quest", name="Slay Dragon", status="active")

# Concepts
kg.add_node("faction_guild", type="faction", name="Thieves Guild", alignment="chaotic")
```

## Edge Structure

Edges connect two nodes with a relationship:

```python
kg.add_edge(
    "player",              # Source node ID
    "merchant",            # Target node ID
    relationship="trusts", # Relationship type
    strength="high",       # Additional data (optional)
    since="3 weeks ago"
)
```

### Common Relationships

```python
# Social relationships
kg.add_edge("char_a", "char_b", relationship="friend")
kg.add_edge("char_a", "char_b", relationship="enemy")
kg.add_edge("child", "parent", relationship="child_of")
kg.add_edge("servant", "master", relationship="serves")

# Location relationships
kg.add_edge("person", "location", relationship="lives_in")
kg.add_edge("person", "location", relationship="visited")

# Ownership
kg.add_edge("person", "item", relationship="owns")
kg.add_edge("person", "item", relationship="wants")

# Quest relationships
kg.add_edge("person", "quest", relationship="assigned")
kg.add_edge("person", "quest", relationship="completed")

# Knowledge
kg.add_edge("person", "secret", relationship="knows_about")
kg.add_edge("person", "event", relationship="witnessed")
```

## Complete Example

Creating a rich NPC memory:

```python
from llm_npc import NPCClient, KnowledgeGraph, Surroundings, Event

# Create knowledge graph
kg = KnowledgeGraph()

# Add characters
kg.add_node("player", type="person", name="Hero", level=10, reputation="hero")
kg.add_node("king", type="person", name="King Aldric", disposition="grateful")
kg.add_node("dragon", type="creature", name="Smaug", status="defeated")
kg.add_node("princess", type="person", name="Princess Elena", status="rescued")

# Add locations
kg.add_node("castle", type="location", name="Royal Castle")
kg.add_node("dragon_lair", type="location", name="Dragon's Lair")

# Add events
kg.add_node("dragon_quest", type="quest", name="Rescue Princess", status="completed")
kg.add_node("battle", type="event", name="Dragon Battle", date="yesterday")

# Add relationships
kg.add_edge("player", "king", relationship="serves")
kg.add_edge("player", "dragon", relationship="defeated")
kg.add_edge("player", "princess", relationship="rescued")
kg.add_edge("player", "dragon_quest", relationship="completed")
kg.add_edge("king", "dragon_quest", relationship="assigned")
kg.add_edge("king", "princess", relationship="father_of")
kg.add_edge("dragon", "princess", relationship="captured")

# Use with NPC
client = NPCClient("http://localhost:8080")

with client.session("game") as session:
    knight = session.create_npc(
        "Sir Reginald",
        "A knight who witnessed the hero's victory over the dragon"
    )
    
    surroundings = Surroundings()
    surroundings.add("Throne Room", "The king's grand throne room")
    surroundings.add("Hero", "The hero who defeated the dragon")
    surroundings.add("King", "King Aldric looking pleased")
    
    events = [
        Event("arrival", "The hero has returned from the dragon quest")
    ]
    
    # NPC has context about all relationships
    response = knight.act(surroundings, events, knowledge_graph=kg)
    
    print(response.text)
    # "The hero has returned victorious! They saved the princess 
    #  and defeated Smaug. The king will surely reward them greatly."
```

## Dynamic Knowledge Graphs

Update knowledge graphs based on game state:

```python
class GameKnowledgeGraph:
    def __init__(self):
        self.kg = KnowledgeGraph()
        self._init_world()
    
    def _init_world(self):
        """Initialize base world knowledge"""
        # Add persistent entities
        self.kg.add_node("player", type="person")
        self.kg.add_node("world", type="meta")
    
    def record_meeting(self, person_a: str, person_b: str):
        """Record that two people met"""
        self.kg.add_edge(person_a, person_b, 
                        relationship="met",
                        timestamp=time.time())
    
    def record_event(self, event_id: str, description: str):
        """Record a world event"""
        self.kg.add_node(event_id, 
                        type="event",
                        description=description,
                        timestamp=time.time())
    
    def set_relationship(self, person_a: str, person_b: str, rel_type: str):
        """Set or update a relationship"""
        self.kg.add_edge(person_a, person_b, relationship=rel_type)
    
    def get_graph(self):
        """Get the current knowledge graph"""
        return self.kg

# Usage
game_kg = GameKnowledgeGraph()

# Record game events
game_kg.record_meeting("player", "merchant")
game_kg.record_event("theft_001", "Theft at the market")
game_kg.set_relationship("player", "merchant", "trusted_by")

# Use with NPC
response = npc.act(surroundings, events, knowledge_graph=game_kg.get_graph())
```

## Knowledge Graph Best Practices

### 1. Use Descriptive IDs

```python
# ✅ Good - clear what it is
kg.add_node("char_merchant_tobias", type="person", name="Tobias")
kg.add_node("loc_market_square", type="location", name="Market Square")

# ❌ Bad - unclear
kg.add_node("n1", type="person", name="Tobias")
kg.add_node("n2", type="location", name="Market Square")
```

### 2. Keep Relevant Information

Only include what the NPC should know:

```python
# ✅ Good - NPC knows about the theft
kg.add_node("theft_01", type="event", description="Theft at market")
kg.add_edge("player", "theft_01", relationship="witnessed")

# ❌ Bad - NPC shouldn't know player's private thoughts
kg.add_node("player_plan", type="thought", description="Player plans to betray king")
```

### 3. Update Relationships

Keep relationships current:

```python
# Initial state
kg.add_edge("guard", "player", relationship="suspicious_of")

# After gaining trust (update or add new edge)
kg.add_edge("guard", "player", relationship="trusts")
```

### 4. Use Meaningful Relationship Types

```python
# ✅ Good - specific and clear
kg.add_edge("player", "sword", relationship="owns")
kg.add_edge("player", "quest", relationship="completed")

# ❌ Bad - too vague
kg.add_edge("player", "sword", relationship="related_to")
kg.add_edge("player", "quest", relationship="has")
```

### 5. Don't Overload

Keep knowledge graphs focused:

```python
# ✅ Good - relevant nodes only (5-15 nodes)
kg.add_node("player", ...)
kg.add_node("king", ...)
kg.add_node("quest", ...)

# ❌ Bad - too much information (50+ nodes)
# This will overwhelm the LLM
for i in range(100):
    kg.add_node(f"node_{i}", ...)
```

Aim for 5-20 nodes and 5-30 edges for best results.

## Knowledge Graph Scope

Different NPCs can have different knowledge graphs:

```python
# Guard's knowledge
guard_kg = KnowledgeGraph()
guard_kg.add_node("player", type="person", reputation="suspicious")
guard_kg.add_edge("player", "theft", relationship="suspect")

guard_response = guard.act([...], knowledge_graph=guard_kg)

# Merchant's knowledge  
merchant_kg = KnowledgeGraph()
merchant_kg.add_node("player", type="person", reputation="good_customer")
merchant_kg.add_edge("player", "merchant", relationship="regular_customer")

merchant_response = merchant.act([...], knowledge_graph=merchant_kg)
```

Different NPCs have different perspectives on the same entity (player).

## Inspecting Knowledge Graphs

```python
kg = KnowledgeGraph()
kg.add_node("a", type="person")
kg.add_node("b", type="person")
kg.add_edge("a", "b", relationship="friend")

# Get as dict
graph_dict = kg.to_dict()
print(graph_dict)
# {
#     'nodes': [{'id': 'a', 'data': {'type': 'person'}}, ...],
#     'edges': [{'source': 'a', 'target': 'b', 'data': {'relationship': 'friend'}}]
# }

# Get size
print(len(kg))  # Total nodes + edges: 3
```

## Next Steps

- **[API Reference: Context](../api/context.md)** - Full KnowledgeGraph API
- **[Examples](../examples/index.md)** - See knowledge graphs in action
- **[Best Practices](../advanced/best-practices.md)** - Advanced patterns and tips

