# LLM NPC Backend Examples

This directory contains example integrations for various game engines and programming languages.

## Available Examples

### [Python](./python/) ✅
Complete game client example demonstrating all backend features.

**Features:**
- Custom tool registration
- NPC management
- Event-driven interactions
- Multi-round inference
- Knowledge graph usage

**Setup:** Uses `uv` for fast dependency management (recommended for macOS)

```bash
cd python
./run.sh
```

### Unity (C#) 🚧
Coming soon - Unity integration example

### Godot (GDScript) 🚧
Coming soon - Godot integration example

### JavaScript/Node.js 🚧
Coming soon - JavaScript client example

## Prerequisites

All examples require:

1. **Backend running** in HTTP mode:
   ```bash
   cd ..
   go build ./cmd/backend/...
   ./backend --http
   ```

2. **(Optional) LLM Provider** for NPC dialogue generation:
   - Ollama: `brew install ollama && ollama serve`
   - LM Studio: Download from [lmstudio.ai](https://lmstudio.ai)

> **Note:** You can test all endpoints except `/npc/act` without an LLM provider.

## Quick Test (No Setup Required)

Test the backend without any example code:

```bash
# Health check
curl http://localhost:8080/health

# Register an NPC
curl -X POST http://localhost:8080/npc/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Test NPC", "background_story": "A test character"}'

# List NPCs
curl http://localhost:8080/npc/list
```

## Example Structure

Each example includes:
- ✅ Complete working code
- ✅ Dependency management
- ✅ README with setup instructions
- ✅ Runnable demo scenario

## Integration Patterns

### Basic Flow
1. Register custom tools (game-specific actions)
2. Register NPCs with background stories
3. Execute NPC actions with surroundings and events
4. Process responses and update game state

### Advanced Features
- **Session Management** - Isolate tools per game instance
- **Knowledge Graphs** - Give NPCs persistent memory
- **Multi-round Inference** - Complex reasoning with tool usage
- **Event Tracking** - NPCs respond to world changes

## Contributing Examples

Want to add an example for your favorite engine?

1. Create a directory: `examples/[language-or-engine]/`
2. Include:
   - Working code
   - Dependency file (requirements.txt, package.json, etc.)
   - README with setup and usage
   - Demo scenario
3. Follow the Python example structure
4. Submit a PR!

## Resources

- [Main Documentation](../README.md)
- [API Reference](../API.md)
- [Testing Guide](../TESTING_GUIDE.md)
- [Documentation Analysis](../DOCUMENTATION_ANALYSIS.md)

## Need Help?

- Check the [troubleshooting section](./python/README.md#troubleshooting) in the Python example
- Review the [API documentation](../API.md)
- Look at the [testing guide](../TESTING_GUIDE.md)

