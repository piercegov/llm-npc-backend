# Simple Game Example

This example demonstrates the core features of the SDK through a tavern scenario with suspicious characters.

## Source Code

The complete example is available at [`examples/simple_game.py`](https://github.com/yourusername/llm-npc-backend/blob/main/python-sdk/examples/simple_game.py).

## Overview

The example creates:

- A tavern innkeeper NPC
- A city guard captain NPC
- Game tools (speak, move, give items)
- A scenario with a suspicious stranger

## Running the Example

```bash
# Start the backend
./backend --http

# Run the example
cd python-sdk/examples
python simple_game.py
```

## Code Walkthrough

### 1. Define Game Tools

```python
from llm_npc import tool

@tool(description="Make the NPC speak dialogue aloud to nearby characters")
def speak(message: str, target: str = None):
    """
    Args:
        message: What the NPC should say
        target: Optional specific character to address
    """
    print(f"[SPEAK] {message}" + (f" (to {target})" if target else ""))

@tool(description="Move the NPC to a different location")
def move_to(location: str):
    """
    Args:
        location: The destination location name
    """
    print(f"[MOVE] Moving to {location}")

@tool(description="Give an item to a character")
def give_item(item: str, recipient: str):
    """
    Args:
        item: The item to give
        recipient: Who to give the item to
    """
    print(f"[GIVE] Giving {item} to {recipient}")
```

### 2. Initialize Client

```python
from llm_npc import NPCClient

client = NPCClient("http://localhost:8080")

# Health check
if not client.health_check():
    print("❌ Backend is not running!")
    return

print("✅ Backend is running")
```

### 3. Create Session and Register Tools

```python
with client.session("simple-game-example") as session:
    # Register tools - SDK automatically converts to backend format
    session.register_tools([speak, move_to, give_item])
    print("✅ Tools registered")
```

### 4. Create NPCs

```python
    innkeeper = session.create_npc(
        name="Elara the Innkeeper",
        background="A warm innkeeper who runs 'The Gilded Swan' tavern. "
                  "She knows all the local gossip and helps travelers."
    )
    
    guard = session.create_npc(
        name="Captain Marcus",
        background="A veteran city guard captain. Fair but strict, "
                  "and currently tracking a group of thieves."
    )
    
    print(f"✅ Created: {innkeeper}")
    print(f"✅ Created: {guard}")
```

### 5. Build Context with Surroundings

```python
from llm_npc import Surroundings, Event

# Build surroundings with simple builder
surroundings = Surroundings()
surroundings.add("Tavern Common Room", 
                "A cozy room with wooden tables, a roaring fireplace, "
                "and the smell of fresh bread. 5 patrons drinking.")
surroundings.add("Hooded Stranger",
                "A figure in a dark cloak sits alone in the corner, "
                "watching the door intently. Hasn't touched their drink.")
surroundings.add("Captain Marcus",
                "The city guard captain just entered through the front door, "
                "scanning the room.")
```

### 6. Define Events

```python
# Define events
events = [
    Event("new_customer", "A hooded stranger entered 10 minutes ago, watching nervously"),
    Event("guard_arrival", "Captain Marcus, the guard captain, just walked in")
]
```

### 7. Make NPC Act

```python
# Execute action - SDK handles all the conversion
response = innkeeper.act(surroundings, events)

if response.success:
    if response.text:
        print(f"💭 Thought: {response.text}")
    
    print(f"Inference rounds: {len(response.rounds)}")
    
    # Check tools used
    if response.tools_used:
        print("🛠️  Tools used:")
        for tool_call in response.tools_used:
            print(f"  - {tool_call.name}({tool_call.args})")
            if tool_call.response:
                print(f"    → {tool_call.response}")
```

### 8. Alternative: Simple Lists

You can also use simple lists for quick prototyping:

```python
# Can also use simple lists (SDK converts automatically)
guard_surroundings = [
    "The Gilded Swan Tavern - a well-maintained inn",
    "Hooded Figure - suspicious person matching thief description",
    "Elara the innkeeper behind the bar",
    "Several regular customers drinking"
]

guard_events = [
    "You spotted someone matching the thieves' description"
]

response = guard.act(guard_surroundings, guard_events)
```

### 9. Using Knowledge Graphs

```python
from llm_npc import KnowledgeGraph

# Build a knowledge graph
kg = KnowledgeGraph()
kg.add_node("stranger_01", type="person", name="Hooded Stranger", suspicious=True)
kg.add_node("guard_marcus", type="person", name="Captain Marcus", role="guard")
kg.add_node("theft_incident", type="event", description="Series of thefts", date="past week")
kg.add_edge("guard_marcus", "theft_incident", relationship="investigating")
kg.add_edge("stranger_01", "theft_incident", relationship="possibly_related")

# New event
new_events = [
    Event("confrontation", "Captain Marcus approached the hooded stranger's table")
]

# NPC with memory context
response = innkeeper.act(surroundings, new_events, knowledge_graph=kg)

if response.success and response.text:
    print(f"💭 Response with memory: {response.text}")
```

## Expected Output

```
=== LLM NPC SDK - Simple Game Example ===

1. Checking backend health...
✅ Backend is running

2. Setting up game session...
✅ Tools registered

3. Creating NPCs...
✅ Created: NPC(id=abc-123, name=Elara the Innkeeper)
✅ Created: NPC(id=def-456, name=Captain Marcus)

4. Running scenario: 'Suspicious Stranger at the Inn'

--- Turn 1: Innkeeper's Perspective ---
💭 Thought: A hooded stranger has been sitting suspiciously in the corner,
and now Captain Marcus has arrived. This could be trouble. I should stay
alert and perhaps inform the captain about the stranger's behavior.

Inference rounds: 1
🛠️  Tools used:
  - speak: {'message': "Captain, there's someone in the corner who's been acting strange.", 'target': 'Captain Marcus'}
[SPEAK] Captain, there's someone in the corner who's been acting strange. (to Captain Marcus)

--- Turn 2: Guard's Perspective ---
💭 Thought: This hooded figure matches the description of the thieves I've
been tracking. I need to approach carefully and question them.

🛠️  Tools used:
  - move_to: {'location': 'Corner table'}
  - speak: {'message': 'Hold there. I need to ask you some questions.', 'target': 'Hooded Figure'}
[MOVE] Moving to Corner table
[SPEAK] Hold there. I need to ask you some questions. (to Hooded Figure)

=== Example Complete ===
```

## Key Takeaways

1. **Simple API**: Just a few lines to create intelligent NPCs
2. **Type Safety**: IDE autocomplete for all methods
3. **Flexible Context**: Use simple lists or builders
4. **Rich Responses**: Easy access to NPC thoughts and actions
5. **Knowledge Graphs**: Give NPCs memory and context

## Modifying the Example

Try changing:

- **Tools**: Add `attack`, `hide`, `search` tools
- **Background**: Give NPCs different personalities
- **Surroundings**: Change the location (forest, castle, dungeon)
- **Events**: Create different scenarios
- **Knowledge Graph**: Add more relationships and history

## Next Steps

- Read the [User Guide](../user-guide/tools.md) for detailed explanations
- Check [API Reference](../api/client.md) for all available methods
- See [Best Practices](../advanced/best-practices.md) for production tips

