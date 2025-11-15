# Client API Reference

The client module provides the main interface for interacting with the LLM NPC Backend.

## NPCClient

::: llm_npc.client.NPCClient
    options:
      show_root_heading: true
      show_source: false
      members:
        - __init__
        - health_check
        - session
        - list_npcs
        - delete_npc

## Session

::: llm_npc.client.Session
    options:
      show_root_heading: true
      show_source: false
      members:
        - __init__
        - register_tools
        - create_npc

## NPC

::: llm_npc.client.NPC
    options:
      show_root_heading: true
      show_source: false
      members:
        - __init__
        - act

## Usage Example

```python
from llm_npc import NPCClient, tool

@tool
def speak(message: str):
    """Make the NPC speak"""
    print(f"NPC says: {message}")

# Create client
client = NPCClient("http://localhost:8080")

# Check health
if client.health_check():
    print("Backend is running")

# Create session
with client.session("game-123") as session:
    # Register tools
    session.register_tools([speak])
    
    # Create NPC
    wizard = session.create_npc(
        name="Merlin",
        background="A wise old wizard"
    )
    
    # Make NPC act
    response = wizard.act(
        surroundings=["Tower room", "Ancient books"],
        events=["A visitor arrives"]
    )
    
    print(response.text)

# List all NPCs
npcs = client.list_npcs()
print(f"Total NPCs: {npcs['count']}")

# Delete an NPC
client.delete_npc(wizard.npc_id)
```

## See Also

- [Getting Started Guide](../getting-started.md)
- [Creating NPCs Guide](../user-guide/npcs.md)
- [Decorators API](decorators.md)

