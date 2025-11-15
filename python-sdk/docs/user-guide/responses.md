# Working with Responses

Every time an NPC acts, you get a `Response` object containing the NPC's thoughts, actions, and metadata. This guide shows you how to work with responses effectively.

## Basic Response Handling

```python
response = npc.act(surroundings, events)

if response.success:
    print(f"NPC thought: {response.text}")
else:
    print(f"Error: {response.error}")
```

## Response Properties

### `success: bool`

Indicates whether the action succeeded:

```python
if response.success:
    # Process the response
    pass
else:
    # Handle error
    print(f"Failed: {response.error}")
```

### `text: str`

The NPC's response or internal thought:

```python
response = npc.act([...])
print(response.text)
# "I should approach the stranger and ask their business here."
```

This is the NPC's natural language response - their thoughts, speech, or narration.

!!! note "Alias Available"
    `response.llm_response` is an alias for `response.text` for compatibility.

### `tools_used: List[ToolCall]`

All tools the NPC used, flattened across all rounds:

```python
for tool in response.tools_used:
    print(f"Tool: {tool.name}")
    print(f"Args: {tool.args}")
    print(f"Success: {tool.success}")
    if tool.result:
        print(f"Result: {tool.result}")
```

Example output:
```
Tool: speak
Args: {'message': 'Halt! State your business.', 'target': 'Stranger'}
Success: True
Result: None
```

### `rounds: List[Round]`

Inference rounds if the NPC performed multi-step reasoning:

```python
print(f"Inference rounds: {len(response.rounds)}")

for i, round in enumerate(response.rounds, 1):
    print(f"\nRound {i}:")
    for tool in round.tools_used:
        print(f"  - {tool.name}: {tool.args}")
```

Most NPCs complete in 1 round, but complex reasoning may use multiple rounds.

### `error: Optional[str]`

Error message if the action failed:

```python
if not response.success:
    print(f"Error occurred: {response.error}")
    # "LLM provider unavailable"
```

### `raw_data: Dict`

Raw backend response for advanced use cases:

```python
raw = response.raw_data
print(raw.keys())
# dict_keys(['success', 'llm_response', 'rounds', 'npc_id', ...])
```

## Tool Calls

Each tool call has these properties:

```python
tool_call = response.tools_used[0]

tool_call.name         # Tool name: "speak"
tool_call.tool_name    # Alias for name
tool_call.args         # Arguments: {"message": "Hello"}
tool_call.success      # Whether it succeeded: True
tool_call.result       # Result/response: "Success"
tool_call.response     # Alias for result
```

### Checking Specific Tools

Check if a specific tool was used:

```python
response = npc.act([...])

# Check if NPC spoke
spoke = any(tool.name == "speak" for tool in response.tools_used)
if spoke:
    print("NPC said something")

# Get all speak tool calls
speech_tools = [t for t in response.tools_used if t.name == "speak"]
for speech in speech_tools:
    print(f"Said: {speech.args['message']}")
```

### Extracting Tool Arguments

```python
response = npc.act([...])

for tool in response.tools_used:
    if tool.name == "move_to":
        destination = tool.args.get("location")
        print(f"NPC moved to: {destination}")
    
    elif tool.name == "attack":
        target = tool.args.get("target")
        weapon = tool.args.get("weapon", "fists")
        print(f"NPC attacked {target} with {weapon}")
```

## Multi-Round Inference

Some NPCs may use multiple reasoning rounds:

```python
response = npc.act([...])

for round_num, round in enumerate(response.rounds, 1):
    print(f"\n=== Round {round_num} ===")
    
    if round.tools_used:
        print("Tools used:")
        for tool in round.tools_used:
            print(f"  - {tool.name}: {tool.args}")
    else:
        print("No tools used (thinking only)")
```

Example output:
```
=== Round 1 ===
Tools used:
  - examine: {'object': 'Suspicious Package'}

=== Round 2 ===
Tools used:
  - alert_guards: {'reason': 'Found explosive device'}
  - move_to: {'location': 'Safe Distance'}
```

## Response Patterns

### Pattern 1: Dialogue NPCs

NPCs that primarily speak:

```python
@tool
def speak(message: str, target: str = None):
    """Speak dialogue"""
    pass

response = npc.act([...])

# Extract all dialogue
for tool in response.tools_used:
    if tool.name == "speak":
        message = tool.args.get("message")
        target = tool.args.get("target", "everyone")
        print(f"To {target}: {message}")
```

### Pattern 2: Action NPCs

NPCs that perform physical actions:

```python
response = npc.act([...])

# Process different action types
for tool in response.tools_used:
    if tool.name == "move_to":
        # Handle movement
        update_npc_position(npc, tool.args["location"])
    
    elif tool.name == "pick_up":
        # Handle item pickup
        add_to_inventory(npc, tool.args["item"])
    
    elif tool.name == "use_item":
        # Handle item usage
        use_item(npc, tool.args["item"])
```

### Pattern 3: Mixed NPCs

NPCs that think, speak, and act:

```python
response = npc.act([...])

# Show thought process
if response.text:
    print(f"💭 {response.text}")

# Show actions
if response.tools_used:
    print("Actions:")
    for tool in response.tools_used:
        if tool.name == "speak":
            print(f"  🗣️  {tool.args['message']}")
        elif tool.name == "move_to":
            print(f"  🚶 Moved to {tool.args['location']}")
        else:
            print(f"  🛠️  {tool.name}: {tool.args}")
```

## Error Handling

Handle different error scenarios:

```python
from llm_npc.exceptions import BackendError

response = npc.act([...])

if not response.success:
    error_msg = response.error or "Unknown error"
    
    if "timeout" in error_msg.lower():
        print("LLM took too long to respond")
    elif "unavailable" in error_msg.lower():
        print("LLM service is down")
    elif "rate limit" in error_msg.lower():
        print("Too many requests")
    else:
        print(f"Error: {error_msg}")
```

## Empty Responses

Sometimes NPCs produce no output:

```python
response = npc.act([...])

if response.success:
    if not response.text and not response.tools_used:
        print("⚠️  NPC produced no output")
        print("   This can happen with small models")
        print("   Try: a larger model, simpler context, or different prompt")
    else:
        # Process response normally
        pass
```

## Caching Responses

Store responses for game state:

```python
class NPCHistory:
    def __init__(self):
        self.history = []
    
    def record(self, npc_name: str, response):
        """Record an NPC's action"""
        self.history.append({
            "npc": npc_name,
            "timestamp": time.time(),
            "thought": response.text,
            "actions": [
                {
                    "tool": tool.name,
                    "args": tool.args
                }
                for tool in response.tools_used
            ]
        })
    
    def get_last_action(self, npc_name: str):
        """Get NPC's last action"""
        for record in reversed(self.history):
            if record["npc"] == npc_name:
                return record
        return None

# Usage
history = NPCHistory()

response = npc.act([...])
history.record("Guard Captain", response)

# Later...
last = history.get_last_action("Guard Captain")
print(f"Last action: {last['actions']}")
```

## Response Validation

Validate response quality:

```python
def validate_response(response) -> bool:
    """Check if response is valid and useful"""
    if not response.success:
        return False
    
    # At least some output
    if not response.text and not response.tools_used:
        return False
    
    # Tools succeeded
    for tool in response.tools_used:
        if not tool.success:
            print(f"Warning: Tool {tool.name} failed")
            return False
    
    return True

response = npc.act([...])
if validate_response(response):
    # Process response
    pass
else:
    # Retry or use fallback
    pass
```

## Combining Multiple NPCs

Process responses from multiple NPCs:

```python
def process_scene(npcs, shared_surroundings, shared_events):
    """Process a scene with multiple NPCs"""
    responses = {}
    
    for npc_name, npc in npcs.items():
        response = npc.act(shared_surroundings, shared_events)
        responses[npc_name] = response
    
    # Compile all dialogue
    all_dialogue = []
    for npc_name, response in responses.items():
        for tool in response.tools_used:
            if tool.name == "speak":
                all_dialogue.append({
                    "speaker": npc_name,
                    "message": tool.args["message"]
                })
    
    return all_dialogue

# Usage
npcs = {
    "Guard": guard_npc,
    "Merchant": merchant_npc,
    "Thief": thief_npc
}

dialogue = process_scene(npcs, surroundings, events)
for line in dialogue:
    print(f"{line['speaker']}: {line['message']}")
```

## Next Steps

- **[Knowledge Graphs](knowledge-graphs.md)** - Give NPCs persistent memory
- **[Examples](../examples/index.md)** - See complete working examples
- **[API Reference: Models](../api/models.md)** - Full Response API documentation

