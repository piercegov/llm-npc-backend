# LLM NPC Python SDK

Welcome to the **LLM NPC Python SDK** documentation! This SDK provides a simple and powerful interface for building intelligent Non-Player Characters (NPCs) using Large Language Models.

## Features

- 🎯 **Simple API**: Clean, intuitive interface for creating and managing NPCs
- 🛠️ **Tool Decorators**: Define game actions with a simple `@tool` decorator
- 📦 **Type-Safe**: Full type hints for better IDE support and error detection
- 🏗️ **Builders**: Convenient builders for surroundings, events, and knowledge graphs
- 🔄 **Flexible**: Support for both simple and advanced usage patterns
- 🧠 **Knowledge Graphs**: Give NPCs persistent memory and understanding of relationships

## Quick Example

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
```

## Installation

Install `uv` (fast Python package installer):

```bash
# macOS and Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
```

Install the SDK from the repository:

```bash
cd python-sdk
uv pip install -e .
```

For development with testing and linting tools:

```bash
uv pip install -e ".[dev]"
```

## Requirements

- Python 3.8+
- Running [LLM NPC Backend](https://github.com/yourusername/llm-npc-backend) server

## Next Steps

<div class="grid cards" markdown>

-   :material-clock-fast:{ .lg .middle } __Getting Started__

    ---

    Install and configure the SDK, create your first NPC

    [:octicons-arrow-right-24: Getting Started](getting-started.md)

-   :material-book-open-variant:{ .lg .middle } __User Guide__

    ---

    Learn how to define tools, build context, and work with NPCs

    [:octicons-arrow-right-24: User Guide](user-guide/tools.md)

-   :material-code-json:{ .lg .middle } __API Reference__

    ---

    Detailed API documentation for all classes and functions

    [:octicons-arrow-right-24: API Reference](api/client.md)

-   :material-application-brackets:{ .lg .middle } __Examples__

    ---

    Complete working examples to get you started quickly

    [:octicons-arrow-right-24: Examples](examples/index.md)

</div>

## Why Use This SDK?

Compare the SDK to manually using the backend API:

=== "With SDK"

    ```python
    from llm_npc import NPCClient, tool, Surroundings
    
    @tool
    def speak(message: str):
        """Make the NPC speak"""
        pass
    
    client = NPCClient("http://localhost:8080")
    
    with client.session("my-game") as session:
        session.register_tools([speak])
        npc = session.create_npc("Gandalf", "A wise wizard")
        
        surroundings = Surroundings()
        surroundings.add("Cave", "A dark cave")
        
        response = npc.act(surroundings)
        print(response.text)
    ```

=== "Without SDK (Manual API)"

    ```python
    import requests
    
    # Register tools
    tools_payload = {
        "session_id": "my-game",
        "tools": [{
            "name": "speak",
            "description": "Make the NPC speak",
            "parameters": {
                "message": {
                    "type": "string",
                    "description": "What to say",
                    "required": True
                }
            }
        }]
    }
    requests.post("http://localhost:8080/tools/register", json=tools_payload)
    
    # Register NPC
    npc_payload = {
        "name": "Gandalf",
        "background_story": "A wise wizard"
    }
    result = requests.post("http://localhost:8080/npc/register", json=npc_payload).json()
    npc_id = result['npc_id']
    
    # Act
    act_payload = {
        "npc_id": npc_id,
        "session_id": "my-game",
        "surroundings": [{"name": "Cave", "description": "A dark cave"}]
    }
    response = requests.post("http://localhost:8080/npc/act", json=act_payload).json()
    print(response.get('llm_response', ''))
    ```

The SDK reduces boilerplate by **~90%** and provides type safety, better error handling, and a more intuitive API!

## Community & Support

- **Repository**: [GitHub](https://github.com/yourusername/llm-npc-backend)
- **Issues**: [Report bugs or request features](https://github.com/yourusername/llm-npc-backend/issues)
- **Backend Docs**: [Main project documentation](https://github.com/yourusername/llm-npc-backend)

## License

This project is licensed under the MIT License.

