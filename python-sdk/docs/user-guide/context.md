# Building Context

NPCs need context to make intelligent decisions. The SDK provides flexible ways to specify what an NPC can see, what events have occurred, and what they know.

## Three Ways to Provide Context

1. **Simple lists** (quick and easy)
2. **Builder classes** (structured and readable)
3. **Context builder** (advanced, complete control)

## Surroundings

Surroundings describe what the NPC can see or interact with in their environment.

### Method 1: Simple Strings

The easiest way - just pass strings:

```python
response = npc.act(
    surroundings=[
        "Dark forest",
        "Rusty sword on ground",
        "Wolf growling nearby"
    ]
)
```

The SDK automatically converts each string to a surrounding object with the string as both name and description.

### Method 2: Dict Format

More control with explicit names and descriptions:

```python
response = npc.act(
    surroundings=[
        {
            "name": "Forest",
            "description": "A dark forest with tall pine trees and thick underbrush"
        },
        {
            "name": "Sword",
            "description": "An old, rusty iron sword lying in the dirt"
        },
        {
            "name": "Wolf",
            "description": "A large grey wolf, teeth bared, growling aggressively"
        }
    ]
)
```

### Method 3: Surroundings Builder (Recommended)

The cleanest approach using the builder:

```python
from llm_npc import Surroundings

surroundings = Surroundings()
surroundings.add("Forest", "A dark forest with tall pine trees")
surroundings.add("Sword", "An old, rusty iron sword lying in the dirt")
surroundings.add("Wolf", "A large grey wolf, teeth bared, growling")

response = npc.act(surroundings)
```

Benefits:

- **Method chaining**: `Surroundings().add(...).add(...)`
- **Type safety**: Better IDE autocomplete
- **Clear structure**: Easy to read and maintain

## Events

Events describe what has happened recently that the NPC should know about.

### Method 1: Simple Strings

Quick event descriptions:

```python
response = npc.act(
    surroundings=[...],
    events=[
        "The wolf started growling",
        "You heard a branch snap behind you"
    ]
)
```

### Method 2: Dict Format

Structured events with types:

```python
response = npc.act(
    surroundings=[...],
    events=[
        {
            "event_type": "threat_detected",
            "event_description": "A wolf appeared and started growling"
        },
        {
            "event_type": "sound",
            "event_description": "A branch snapped behind you"
        }
    ]
)
```

### Method 3: Event Objects

Using the `Event` class:

```python
from llm_npc import Event

events = [
    Event("threat_detected", "A wolf appeared and started growling"),
    Event("sound", "A branch snapped behind you")
]

response = npc.act(surroundings, events)
```

## Complete Context Example

Combining surroundings and events:

```python
from llm_npc import Surroundings, Event

# Build surroundings
surroundings = Surroundings()
surroundings.add("Tavern", "A crowded tavern with wooden tables and a fireplace")
surroundings.add("Bartender", "A gruff bartender cleaning glasses")
surroundings.add("Stranger", "A hooded figure in the corner watching you")

# Define events
events = [
    Event("arrival", "You just entered the tavern"),
    Event("observation", "The hooded stranger has been watching you for 5 minutes")
]

# Make NPC act with full context
response = npc.act(surroundings, events)
```

## Advanced: Context Builder

For complex scenarios, use `ContextBuilder`:

```python
from llm_npc import ContextBuilder, KnowledgeGraph

# Build complete context
context = ContextBuilder()

# Add surroundings
context.add_surrounding("Castle Throne Room", "A grand hall with marble floors")
context.add_surrounding("King", "An elderly king sitting on an ornate throne")
context.add_surrounding("Guards", "Two armored guards flanking the throne")

# Add events
context.add_event("summons", "You were summoned by the king")
context.add_event("quest_completed", "You completed the dragon quest")

# Add knowledge graph (optional)
kg = KnowledgeGraph()
kg.add_node("king", type="person", disposition="grateful")
kg.add_node("dragon_quest", type="quest", status="completed")
kg.add_edge("king", "dragon_quest", relationship="assigned")
context.set_knowledge_graph(kg)

# Use the complete context
response = npc.act(context)
```

The context builder creates a complete, structured context with surroundings, events, and optional knowledge graph.

## Best Practices

### 1. Be Specific

Detailed descriptions help NPCs make better decisions:

```python
# ❌ Vague
surroundings.add("Room", "A room")

# ✅ Specific
surroundings.add("Throne Room", 
    "A vast throne room with vaulted ceilings, stained glass windows, "
    "and red carpet leading to an ornate golden throne")
```

### 2. Include Relevant Details

Add details that might affect NPC behavior:

```python
# Characters
surroundings.add("Guard", 
    "A heavily armored guard, hand on sword hilt, watching you suspiciously")

# Objects
surroundings.add("Chest", 
    "An unlocked wooden chest with the lid slightly ajar, smelling of old leather")

# Environment
surroundings.add("Weather", 
    "Heavy rain pounding on the windows, thunder rumbling in the distance")
```

### 3. Temporal Events

Events should explain **what just happened** or **what's currently happening**:

```python
events = [
    Event("action", "The guard just drew his sword"),
    Event("dialogue", "The merchant shouted 'Thief!'"),
    Event("change", "The door slammed shut and locked"),
]
```

### 4. Scope Appropriately

Only include what the NPC can perceive:

```python
# ✅ Good - NPC can see/hear these
surroundings = Surroundings()
surroundings.add("Market", "Busy marketplace with vendor stalls")
surroundings.add("Crowd", "Dozens of shoppers browsing goods")

# ❌ Bad - NPC can't know this
surroundings.add("Secret Plot", "The king is secretly planning your arrest")

# ✅ Better - observable signs
surroundings.add("Guards", "More guards than usual, watching the crowd closely")
events = [Event("observation", "Guards seem alert and searching for someone")]
```

### 5. Reusable Context

Build reusable context templates:

```python
def create_market_scene():
    """Standard market scene"""
    surroundings = Surroundings()
    surroundings.add("Market Square", "A busy open-air market")
    surroundings.add("Vendors", "Various vendors selling goods")
    surroundings.add("Shoppers", "Crowd of people shopping")
    return surroundings

# Use for multiple NPCs
merchant_response = merchant.act(create_market_scene())
guard_response = guard.act(create_market_scene())
```

## Dynamic Context

Update context based on game state:

```python
def get_current_scene(time_of_day, weather, npcs_present):
    surroundings = Surroundings()
    
    # Base location
    surroundings.add("Town Square", "The central square of the town")
    
    # Time-based
    if time_of_day == "night":
        surroundings.add("Lighting", "Dim light from scattered torches")
    else:
        surroundings.add("Lighting", "Bright sunlight")
    
    # Weather
    if weather == "rain":
        surroundings.add("Weather", "Heavy rain, puddles forming")
    
    # NPCs
    for npc_name in npcs_present:
        surroundings.add(npc_name, f"{npc_name} is here")
    
    return surroundings

# Use dynamic context
scene = get_current_scene("night", "rain", ["Guard", "Merchant"])
response = npc.act(scene)
```

## Context Limits

Be mindful of context size:

```python
# ✅ Good - focused context
surroundings = Surroundings()
surroundings.add("Room", "...")      # 3-4 key items
surroundings.add("Person", "...")
surroundings.add("Object", "...")

# ⚠️ Too much - may overwhelm the LLM
surroundings = Surroundings()
for i in range(50):  # Don't do this!
    surroundings.add(f"Item {i}", "...")
```

Aim for 3-10 key surroundings and 1-5 events per action.

## Iterating Context

Access surroundings programmatically:

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
```

## Next Steps

- **[Working with Responses](responses.md)** - Parse NPC actions and responses
- **[Knowledge Graphs](knowledge-graphs.md)** - Give NPCs long-term memory
- **[API Reference: Context](../api/context.md)** - Full context API documentation

