# Python Examples for LLM NPC Backend

This directory contains Python examples demonstrating how to integrate with the LLM NPC Backend.

## Prerequisites

1. **uv** - Fast Python package installer (recommended for macOS with Homebrew)
   ```bash
   brew install uv
   ```

2. **Backend running** - Start the backend in HTTP mode:
   ```bash
   cd ../..  # Go to project root
   go build ./cmd/backend/...
   ./backend --http
   ```

3. **(Optional) LLM Provider** - For NPC actions to work:
   - **Ollama**: `brew install ollama && ollama serve`
     - **Recommended models:**
       - `ollama pull llama3:8b` (best quality)
       - `ollama pull mistral:7b` (good balance)
       - `ollama pull qwen3:1.7b` (fast, but may produce empty responses)
   - **LM Studio**: Download and run from [lmstudio.ai](https://lmstudio.ai)

## Quick Start

### Option 1: Using the run script (recommended)

```bash
./run.sh
```

The script automatically:
- Checks if `uv` is installed
- Creates a virtual environment
- Installs dependencies
- Runs the example

### Option 2: Manual with uv

```bash
# Install dependencies
uv pip install -r <(uv pip compile pyproject.toml)

# Run the example
uv run game_client.py
```

### Option 3: Traditional venv

```bash
# Create virtual environment
python3 -m venv .venv
source .venv/bin/activate

# Install dependencies
pip install -e .

# Run the example
python game_client.py
```

## What the Example Does

The `game_client.py` demonstrates a complete game integration:

1. **Health Check** - Verifies backend is running
2. **Tool Registration** - Registers custom game tools:
   - `speak` - Make NPCs speak dialogue
   - `move_to` - Move NPCs to locations
   - `give_item` - Transfer items between characters
3. **NPC Registration** - Creates two NPCs:
   - Elara the Innkeeper
   - Captain Marcus (city guard)
4. **Game Scenario** - Simulates "Suspicious Stranger at the Inn"
   - Multiple NPC perspectives
   - Event-driven interactions
   - Multi-round inference
5. **Knowledge Graph** - Demonstrates NPC memory with relationships

## Expected Output

```
=== LLM NPC Backend - Example Game Client ===

1. Checking backend health...
✓ Backend is running

2. Registering custom game tools...
✓ Registered 3 tools: ['speak', 'move_to', 'give_item']

3. Registering NPCs...
✓ Registered NPC 'Elara the Innkeeper' with ID: [uuid]
✓ Registered NPC 'Captain Marcus' with ID: [uuid]

...
```

## Modifying the Example

### Change the Backend URL

Edit `game_client.py`:

```python
# For Unix socket mode
client = NPCBackendClient(base_url="http://unix")

# For different HTTP port
client = NPCBackendClient(base_url="http://localhost:3000")
```

### Add Your Own Tools

```python
my_tools = [
    {
        "name": "attack",
        "description": "Attack another character",
        "parameters": {
            "target": {
                "type": "string",
                "description": "Who to attack",
                "required": True
            }
        }
    }
]

client.register_tools("my-session", my_tools)
```

### Create Your Own NPCs

```python
merchant_id = client.register_npc(
    name="Grubak the Merchant",
    background_story="A cunning goblin merchant who trades in exotic goods..."
)
```

## Troubleshooting

### "Backend is not running"
- Make sure you've started the backend: `./backend --http`
- Check it's accessible: `curl http://localhost:8080/health`

### "LLM service is currently unavailable"
- This is expected if you don't have Ollama or LM Studio running
- The example will still work for registration, listing, etc.
- To test NPC actions, start an LLM provider

### Empty NPC responses
- **Symptom**: NPCs return empty text, or only use tools without dialogue
- **Cause**: Small language models (< 3B parameters) may struggle with complex prompts
- **Solution**: Use a larger model:
  ```bash
  ollama pull llama3:8b        # Recommended
  # or
  ollama pull mistral:7b        # Good alternative
  ```
- **Note**: Tool-only responses are valid - NPCs can act without speaking!
- To change models, set in `.env`: `OLLAMA_MODEL=llama3:8b`

### "Module not found: requests"
- Make sure you ran `uv` or installed dependencies
- If using manual venv, activate it: `source .venv/bin/activate`

### uv command not found
- Install uv: `brew install uv`
- Or use traditional Python: `python3 -m venv .venv && source .venv/bin/activate && pip install -e .`

## Next Steps

1. Modify the example scenario
2. Add your own NPCs and tools
3. Integrate into your game engine
4. Check out the other examples (coming soon):
   - Unity C# example
   - Godot GDScript example
   - JavaScript/Node.js example

## Resources

- [Main Project README](../../README.md)
- [API Documentation](../../API.md)
- [Testing Guide](../../TESTING_GUIDE.md)
- [uv Documentation](https://github.com/astral-sh/uv)

