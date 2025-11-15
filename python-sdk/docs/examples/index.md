# Examples

Complete working examples to help you get started with the SDK.

## Available Examples

### [Simple Game Example](simple-game.md)

A complete tavern scenario demonstrating:

- Defining tools with `@tool` decorator
- Creating multiple NPCs
- Using surroundings and events
- Knowledge graphs for NPC memory
- Processing NPC responses

**Best for**: Understanding the basics and seeing the SDK in action

## Running the Examples

All examples are in the `examples/` directory:

```bash
cd python-sdk/examples
python simple_game.py
```

!!! note "Backend Required"
    Make sure the backend is running before executing examples:
    ```bash
    # In repository root
    ./backend --http
    ```

## Example Patterns

### Quick Start Pattern

Minimal code to get an NPC working:

```python
from llm_npc import NPCClient, tool

@tool
def speak(message: str):
    """Speak dialogue"""
    pass

client = NPCClient("http://localhost:8080")

with client.session("game") as session:
    session.register_tools([speak])
    npc = session.create_npc("Guard", "A city guard")
    response = npc.act(["Town square"], ["A stranger approaches"])
    print(response.text)
```

### Multiple NPCs Pattern

Managing several NPCs:

```python
from llm_npc import NPCClient, tool, Surroundings

@tool
def speak(message: str):
    pass

client = NPCClient("http://localhost:8080")

with client.session("game") as session:
    session.register_tools([speak])
    
    # Create NPCs
    guard = session.create_npc("Guard", "Stern city guard")
    merchant = session.create_npc("Merchant", "Friendly merchant")
    thief = session.create_npc("Thief", "Sneaky thief")
    
    # Shared surroundings
    surroundings = Surroundings()
    surroundings.add("Market", "Busy marketplace")
    
    # Each NPC responds
    guard_response = guard.act(surroundings, ["Watching for trouble"])
    merchant_response = merchant.act(surroundings, ["Selling goods"])
    thief_response = thief.act(surroundings, ["Looking for opportunities"])
```

### State Management Pattern

Tracking NPC actions over time:

```python
class GameState:
    def __init__(self):
        self.npc_locations = {}
        self.npc_inventory = {}
    
    def process_response(self, npc_name, response):
        """Process NPC action and update state"""
        for tool in response.tools_used:
            if tool.name == "move_to":
                self.npc_locations[npc_name] = tool.args["location"]
            elif tool.name == "pick_up":
                if npc_name not in self.npc_inventory:
                    self.npc_inventory[npc_name] = []
                self.npc_inventory[npc_name].append(tool.args["item"])

# Usage
state = GameState()
response = npc.act(surroundings)
state.process_response("Guard", response)
```

### Dynamic Context Pattern

Building context from game state:

```python
def build_context(player_location, time_of_day, weather):
    """Build context dynamically"""
    from llm_npc import Surroundings, Event
    
    surroundings = Surroundings()
    surroundings.add("Location", player_location)
    
    # Add time-based details
    if time_of_day == "night":
        surroundings.add("Lighting", "Dark, moonlit")
    else:
        surroundings.add("Lighting", "Bright daylight")
    
    # Add weather
    if weather == "rain":
        events = [Event("weather", "It starts raining heavily")]
    else:
        events = []
    
    return surroundings, events

# Use
surroundings, events = build_context("Forest", "night", "rain")
response = npc.act(surroundings, events)
```

### Error Handling Pattern

Robust error handling:

```python
from llm_npc import NPCClient
from llm_npc.exceptions import BackendError, LLMNPCError

def safe_npc_act(npc, surroundings, events, max_retries=3):
    """Make NPC act with retry logic"""
    for attempt in range(max_retries):
        try:
            response = npc.act(surroundings, events)
            if response.success:
                return response
            else:
                print(f"Attempt {attempt + 1} failed: {response.error}")
        except BackendError as e:
            if "timeout" in str(e).lower() and attempt < max_retries - 1:
                print(f"Timeout, retrying... ({attempt + 1}/{max_retries})")
                continue
            raise
    
    return None

# Usage
response = safe_npc_act(npc, surroundings, events)
if response:
    print(response.text)
else:
    print("All attempts failed")
```

## More Examples

Check the `examples/` directory for additional examples and patterns. Contributions welcome!

## Next Steps

- Try modifying the [Simple Game Example](simple-game.md)
- Check the [User Guide](../user-guide/tools.md) for detailed explanations
- See [Best Practices](../advanced/best-practices.md) for production tips

