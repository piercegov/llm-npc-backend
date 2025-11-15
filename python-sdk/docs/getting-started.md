# Getting Started

This guide will help you get up and running with the LLM NPC Python SDK in minutes.

## Prerequisites

Before you begin, ensure you have:

1. **Python 3.8 or higher** installed
2. **LLM NPC Backend** running (see [backend setup](https://github.com/yourusername/llm-npc-backend#readme))
3. **pip** package manager

## Installation

### Install uv

First, install `uv` - a fast Python package installer:

```bash
# macOS and Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"

# Or with pip
pip install uv
```

### From Source

Clone the repository and install the SDK:

```bash
# Clone the repository
git clone https://github.com/yourusername/llm-npc-backend.git
cd llm-npc-backend/python-sdk

# Install in development mode
uv pip install -e .
```

### With Development Tools

If you're planning to contribute or modify the SDK:

```bash
uv pip install -e ".[dev]"
```

This installs additional tools like pytest, black, and mypy.

## Starting the Backend

The SDK communicates with the LLM NPC Backend server. Make sure it's running:

```bash
# In the repository root
./backend --http

# Or build and run
go build ./cmd/backend/...
./backend --http
```

The server will start on `http://localhost:8080` by default.

!!! tip "Health Check"
    You can verify the backend is running:
    ```bash
    curl http://localhost:8080/health
    # Should return: pong
    ```

## Your First NPC

Let's create a simple NPC that can speak and move.

### Step 1: Define Tools

Tools are actions that NPCs can perform. Use the `@tool` decorator:

```python
from llm_npc import tool

@tool(description="Make the NPC speak dialogue")
def speak(message: str, target: str = None):
    """
    Args:
        message: What the NPC should say
        target: Optional character to address
    """
    print(f"🗣️  {message}")
    if target:
        print(f"   (to {target})")

@tool(description="Move to a new location")
def move_to(location: str):
    """
    Args:
        location: Where to move
    """
    print(f"🚶 Moving to {location}")
```

### Step 2: Create a Client and Session

```python
from llm_npc import NPCClient

# Initialize the client
client = NPCClient("http://localhost:8080")

# Check if backend is running
if not client.health_check():
    print("❌ Backend not running!")
    exit(1)

print("✅ Connected to backend")
```

### Step 3: Register Tools and Create an NPC

```python
# Create a session for your game
with client.session("my-first-game") as session:
    # Register your tools
    session.register_tools([speak, move_to])
    print("✅ Tools registered")
    
    # Create an NPC
    guard = session.create_npc(
        name="Sir Reginald",
        background="A loyal knight who guards the castle gate. "
                   "He's friendly to visitors but wary of suspicious characters."
    )
    print(f"✅ Created NPC: {guard}")
```

### Step 4: Make the NPC Act

```python
    # Define what the NPC can see/experience
    response = guard.act(
        surroundings=[
            "Castle Gate - A large wooden gate with iron reinforcements",
            "Traveler - A hooded figure approaching the gate",
            "Sunset - The sun is setting behind the mountains"
        ],
        events=[
            "A hooded traveler is approaching the gate"
        ]
    )
    
    # Check the response
    if response.success:
        print(f"\n💭 NPC Thought: {response.text}")
        
        # See what actions the NPC took
        if response.tools_used:
            print("\n🛠️  Actions taken:")
            for tool in response.tools_used:
                print(f"  - {tool.name}: {tool.args}")
    else:
        print(f"❌ Error: {response.error}")
```

### Complete Example

Here's the complete code:

```python
from llm_npc import NPCClient, tool

# Define tools
@tool(description="Make the NPC speak dialogue")
def speak(message: str, target: str = None):
    """
    Args:
        message: What the NPC should say
        target: Optional character to address
    """
    print(f"🗣️  {message}")

@tool(description="Move to a new location")
def move_to(location: str):
    """
    Args:
        location: Where to move
    """
    print(f"🚶 Moving to {location}")

# Initialize client
client = NPCClient("http://localhost:8080")

if not client.health_check():
    print("❌ Backend not running!")
    exit(1)

# Create session and NPC
with client.session("my-first-game") as session:
    session.register_tools([speak, move_to])
    
    guard = session.create_npc(
        name="Sir Reginald",
        background="A loyal knight who guards the castle gate"
    )
    
    # Make the NPC act
    response = guard.act(
        surroundings=[
            "Castle Gate",
            "Traveler approaching",
            "Sunset"
        ],
        events=["A hooded traveler approaches"]
    )
    
    if response.success:
        print(f"💭 {response.text}")
        for tool in response.tools_used:
            print(f"🛠️  {tool.name}: {tool.args}")
```

## Running the Example

Save the code to `first_npc.py` and run:

```bash
python first_npc.py
```

You should see output like:

```
✅ Connected to backend
✅ Tools registered
✅ Created NPC: NPC(id=abc-123, name=Sir Reginald)
💭 A traveler is approaching. I should greet them and inquire about their business.
🛠️  speak: {'message': 'Halt! State your business at the castle.', 'target': 'Traveler'}
🗣️  Halt! State your business at the castle.
   (to Traveler)
```

## What's Next?

Now that you have your first NPC working, explore more advanced features:

- **[Defining Tools](user-guide/tools.md)** - Learn about parameter types, docstrings, and advanced tool patterns
- **[Building Context](user-guide/context.md)** - Use builders for more structured surroundings and events
- **[Knowledge Graphs](user-guide/knowledge-graphs.md)** - Give your NPCs memory and understanding of relationships
- **[Working with Responses](user-guide/responses.md)** - Parse NPC responses and handle multi-round inference

## Troubleshooting

### Backend Connection Issues

If you get connection errors:

```python
from llm_npc.exceptions import BackendConnectionError

try:
    client.health_check()
except BackendConnectionError:
    print("Cannot connect to backend. Is it running?")
```

### Empty Responses

If your NPC produces empty responses:

- **Try a larger model**: Small models like `qwen3:1.7b` may not perform well
- **Check backend logs**: Look for LLM errors
- **Simplify the scenario**: Start with simpler surroundings

```python
response = npc.act(["Simple room", "Door"])
```

### Tool Registration Errors

Make sure your functions are decorated:

```python
# ✅ Correct
@tool
def my_action():
    pass

# ❌ Wrong - missing decorator
def my_action():
    pass
```

## Configuration

### Custom Backend URL

```python
client = NPCClient("http://your-server:8080")
```

### Using pip Instead of uv

If you prefer using `pip`, simply replace `uv pip` with `pip`:

```bash
pip install -e .
pip install -e ".[dev]"
```

### Environment Variables

The backend can be configured via environment variables. See the [backend documentation](https://github.com/yourusername/llm-npc-backend#configuration) for details.

## Next Steps

Ready to dive deeper? Check out the **[User Guide](user-guide/tools.md)** for comprehensive tutorials on all SDK features!

