# LLM NPC Python SDK

A simple and powerful Python SDK for building intelligent NPCs using the [LLM NPC Backend](../README.md).

## Features

- 🎯 **Simple API**: Clean, intuitive interface for creating and managing NPCs
- 🛠️ **Tool Decorators**: Define game actions with a simple `@tool` decorator
- 📦 **Type-Safe**: Full type hints for better IDE support
- 🏗️ **Builders**: Convenient builders for surroundings, events, and knowledge graphs
- 🔄 **Flexible**: Support for both simple and advanced usage patterns

## Installation

### From Source (Development)

```bash
# From the repository root
cd python-sdk
pip install -e .
```

### With Development Dependencies

```bash
pip install -e ".[dev]"
```

## Quick Start

```python
from llm_npc import NPCClient, tool

# Define game tools with decorators
@tool(description="Make the NPC speak")
def speak(message: str, target: str = None):
    """
    Args:
        message: What the NPC should say
        target: Optional character to address
    """
    print(f"NPC says: {message}")

# Initialize client
client = NPCClient("http://localhost:8080")

# Create a session and NPC
with client.session("my-game-session") as session:
    # Register tools
    session.register_tools([speak])
    
    # Create an NPC
    wizard = session.create_npc(
        name="Gandalf",
        background="A wise wizard with centuries of experience"
    )
    
    # Make the NPC act
    response = wizard.act(
        surroundings=["Dark cave", "Sleeping dragon", "Pile of gold"],
        events=["You hear footsteps behind you"]
    )
    
    print(response.text)
    for tool in response.tools_used:
        print(f"Used tool: {tool.name}")
```

## Usage Guide

### 1. Defining Tools

Tools are game-specific actions that NPCs can perform. Define them using the `@tool` decorator:

```python
from llm_npc import tool

@tool
def attack(target: str, weapon: str = "fists"):
    """
    Attack a target.
    
    Args:
        target: The enemy to attack
        weapon: Weapon to use (default: fists)
    """
    return f"Attacked {target} with {weapon}"

@tool(description="Cast a magical spell")
def cast_spell(spell_name: str, target: str):
    """
    Args:
        spell_name: Name of the spell to cast
        target: Who/what to cast it on
    """
    return f"Cast {spell_name} on {target}"
```

The SDK automatically:
- Extracts the function name as the tool name
- Uses the docstring as the tool description
- Parses parameters from the function signature
- Converts type hints to the backend schema format

### 2. Creating NPCs

```python
from llm_npc import NPCClient

client = NPCClient("http://localhost:8080")

with client.session("game-session-123") as session:
    # Register your tools
    session.register_tools([attack, cast_spell])
    
    # Create NPCs
    warrior = session.create_npc(
        name="Conan",
        background="A mighty barbarian warrior from the north"
    )
    
    mage = session.create_npc(
        name="Merlin",
        background="An ancient wizard with vast magical knowledge"
    )
```

### 3. Building Context

#### Simple Lists (Automatic Conversion)

```python
# Simple strings - SDK converts automatically
response = warrior.act(
    surroundings=["Forest clearing", "Rusty sword on ground", "Angry wolf"],
    events=["The wolf growls at you"]
)
```

#### Using Builders (Recommended)

```python
from llm_npc import Surroundings, Event

# Build surroundings
surroundings = Surroundings()
surroundings.add("Castle Throne Room", "A grand hall with tall pillars")
surroundings.add("King", "An elderly king sitting on his throne")
surroundings.add("Guards", "Two armored guards flanking the throne")

# Create events
events = [
    Event("arrival", "You were summoned by the king"),
    Event("quest_offered", "The king mentions a quest")
]

response = warrior.act(surroundings, events)
```

#### Advanced: Context Builder

```python
from llm_npc import ContextBuilder, KnowledgeGraph

# Build complex context
context = ContextBuilder()
context.add_surrounding("Village Square", "A busy marketplace")
context.add_surrounding("Merchant", "A nervous-looking merchant")
context.add_event("theft", "You witnessed a pickpocket")

# Add knowledge graph (NPC memory)
kg = KnowledgeGraph()
kg.add_node("merchant_bob", type="person", role="merchant", trust_level="high")
kg.add_node("thief_guild", type="organization", reputation="notorious")
kg.add_edge("merchant_bob", "thief_guild", relationship="victim_of")

context.set_knowledge_graph(kg)

# Use the complete context
response = warrior.act(context)
```

### 4. Working with Responses

```python
response = wizard.act(surroundings, events)

# Check if successful
if response.success:
    # Get the NPC's thought/response
    print(f"NPC thought: {response.text}")
    
    # Check how many inference rounds
    print(f"Rounds: {len(response.rounds)}")
    
    # See what tools were used
    for tool_call in response.tools_used:
        print(f"Tool: {tool_call.name}")
        print(f"Args: {tool_call.args}")
        print(f"Success: {tool_call.success}")
        print(f"Result: {tool_call.result}")
    
    # Access raw backend data if needed
    raw_data = response.raw_data
else:
    print(f"Error: {response.error}")
```

### 5. Knowledge Graphs (NPC Memory)

Knowledge graphs let you give NPCs persistent memory and understanding of relationships:

```python
from llm_npc import KnowledgeGraph

kg = KnowledgeGraph()

# Add entities (nodes)
kg.add_node("player", type="person", name="Hero", reputation="good")
kg.add_node("dragon", type="creature", name="Smaug", threat_level="extreme")
kg.add_node("treasure", type="object", value="immense")

# Add relationships (edges)
kg.add_edge("dragon", "treasure", relationship="guards")
kg.add_edge("player", "dragon", relationship="hunting")

# Use in NPC action
response = npc.act(surroundings, events, knowledge_graph=kg)
```

### 6. Session Management

Sessions group tools and NPCs together:

```python
client = NPCClient("http://localhost:8080")

# Context manager (recommended)
with client.session("my-session") as session:
    session.register_tools([...])
    npc = session.create_npc(...)
    # Session automatically cleaned up

# Manual management
session = client.session("my-session")
session.register_tools([...])
npc = session.create_npc(...)
```

### 7. Error Handling

```python
from llm_npc import (
    BackendConnectionError,
    BackendError,
    ToolRegistrationError
)

try:
    response = npc.act(surroundings)
except BackendConnectionError:
    print("Cannot connect to backend - is it running?")
except BackendError as e:
    print(f"Backend error: {e}")
    print(f"Status code: {e.status_code}")
except ToolRegistrationError:
    print("Failed to register tools")
```

## Examples

See the [examples](./examples/) directory for complete working examples:

- [`simple_game.py`](./examples/simple_game.py) - Demonstrates the simplified API

Compare this to the [manual API usage](../examples/python/game_client.py) to see how much simpler the SDK is!

## API Reference

### NPCClient

```python
NPCClient(base_url: str = "http://localhost:8080")
```

Main client for interacting with the backend.

Methods:
- `health_check() -> bool` - Check if backend is running
- `session(session_id: str) -> Session` - Create a new session
- `list_npcs() -> Dict` - List all registered NPCs
- `delete_npc(npc_id: str) -> bool` - Delete an NPC

### Session

Manages tools and NPCs for a game session.

Methods:
- `register_tools(tools: List[Callable]) -> Session` - Register tool functions
- `create_npc(name: str, background: str) -> NPC` - Create a new NPC

### NPC

Represents an intelligent NPC.

Methods:
- `act(surroundings, events=None, knowledge_graph=None) -> Response` - Execute an action

### Decorators

- `@tool` - Mark a function as a tool NPCs can use
- `@tool(description="...")` - Provide custom description

### Builders

- `Surroundings()` - Build surroundings list
- `Event(type, description)` - Create an event
- `ContextBuilder()` - Build complete context
- `KnowledgeGraph()` - Build knowledge graphs

### Models

- `Response` - NPC action response
  - `success: bool`
  - `text: str` - NPC's thought/response
  - `rounds: List[Round]` - Inference rounds
  - `tools_used: List[ToolCall]` - All tools used
  - `error: Optional[str]`

- `ToolCall` - Represents a tool usage
  - `name: str`
  - `args: Dict`
  - `success: bool`
  - `result: Optional[str]`

## Development

### Running Tests

```bash
pytest
```

### Code Formatting

```bash
black llm_npc/
```

### Type Checking

```bash
mypy llm_npc/
```

## Requirements

- Python 3.8+
- `requests` library
- Running LLM NPC Backend server

## License

MIT

## Contributing

Contributions are welcome! This SDK is part of the [LLM NPC Backend](../) project.

