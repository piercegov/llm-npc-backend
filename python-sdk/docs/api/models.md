# Models API Reference

The models module provides typed data classes for working with NPC responses and context.

## Response

::: llm_npc.models.Response
    options:
      show_root_heading: true
      show_source: false

## ToolCall

::: llm_npc.models.ToolCall
    options:
      show_root_heading: true
      show_source: false

## Round

::: llm_npc.models.Round
    options:
      show_root_heading: true
      show_source: false

## Surrounding

::: llm_npc.models.Surrounding
    options:
      show_root_heading: true
      show_source: false

## Event

::: llm_npc.models.Event
    options:
      show_root_heading: true
      show_source: false

## Usage Examples

### Working with Responses

```python
response = npc.act(surroundings, events)

# Check success
if response.success:
    # Get NPC's thought
    print(f"Thought: {response.text}")
    
    # Check inference rounds
    print(f"Rounds: {len(response.rounds)}")
    
    # Get all tools used
    for tool in response.tools_used:
        print(f"Tool: {tool.name}")
        print(f"Args: {tool.args}")
        print(f"Success: {tool.success}")
else:
    print(f"Error: {response.error}")
```

### Creating Typed Surroundings

```python
from llm_npc.models import Surrounding

surroundings = [
    Surrounding("Forest", "A dark forest"),
    Surrounding("Sword", "A rusty sword on the ground")
]

response = npc.act(surroundings)
```

### Creating Typed Events

```python
from llm_npc.models import Event

events = [
    Event("arrival", "A stranger approaches"),
    Event("threat", "The stranger draws a weapon")
]

response = npc.act(surroundings, events)
```

### Accessing Tool Calls

```python
response = npc.act(surroundings, events)

for tool_call in response.tools_used:
    # Access properties
    tool_name = tool_call.name  # or tool_call.tool_name
    arguments = tool_call.args
    succeeded = tool_call.success
    result = tool_call.result  # or tool_call.response
    
    # Process based on tool
    if tool_name == "speak":
        message = arguments.get("message")
        print(f"NPC said: {message}")
```

### Iterating Rounds

```python
response = npc.act(surroundings, events)

for round_num, round_obj in enumerate(response.rounds, 1):
    print(f"Round {round_num}:")
    
    for tool in round_obj.tools_used:
        print(f"  - {tool.name}: {tool.args}")
    
    # Access raw data if needed
    raw = round_obj.raw_data
```

## See Also

- [Working with Responses Guide](../user-guide/responses.md)
- [Building Context Guide](../user-guide/context.md)
- [Context API](context.md)

