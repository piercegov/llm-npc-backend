# Defining Tools

Tools are game-specific actions that NPCs can perform. The SDK makes it easy to define tools using the `@tool` decorator, which automatically extracts metadata from your Python functions.

## Basic Tool Definition

The simplest way to define a tool:

```python
from llm_npc import tool

@tool
def speak(message: str):
    """Make the NPC speak dialogue"""
    print(f"NPC says: {message}")
```

The SDK automatically:

- Uses the function name as the tool name (`speak`)
- Extracts the description from the docstring
- Determines parameter types from type hints
- Identifies required vs optional parameters

## Custom Description

Override the docstring description:

```python
@tool(description="Make the NPC speak dialogue aloud to nearby characters")
def speak(message: str):
    """This docstring is ignored when description is provided"""
    print(f"NPC says: {message}")
```

## Parameter Types

The SDK supports multiple parameter types:

```python
@tool
def attack(
    target: str,           # String parameter
    damage: int = 10,      # Integer with default
    critical: bool = False, # Boolean
    multiplier: float = 1.0 # Float
):
    """
    Attack a target.
    
    Args:
        target: Enemy to attack
        damage: Damage amount
        critical: Whether it's a critical hit
        multiplier: Damage multiplier
    """
    actual_damage = damage * multiplier * (2 if critical else 1)
    print(f"Attacked {target} for {actual_damage} damage")
```

### Supported Types

| Python Type | Schema Type | Example |
|-------------|-------------|---------|
| `str` | `string` | `"hello"` |
| `int` | `integer` | `42` |
| `float` | `float` | `3.14` |
| `bool` | `boolean` | `True` |

## Optional Parameters

Parameters with default values are optional:

```python
@tool
def cast_spell(
    spell_name: str,        # Required
    target: str = None,     # Optional
    power: int = 5          # Optional with default
):
    """
    Cast a magical spell.
    
    Args:
        spell_name: Name of the spell to cast
        target: Optional target for the spell
        power: Spell power level
    """
    if target:
        print(f"Cast {spell_name} on {target} (power: {power})")
    else:
        print(f"Cast {spell_name} (power: {power})")
```

## Documenting Parameters

Use Google-style docstrings to document parameters:

```python
@tool
def give_item(item: str, recipient: str, quantity: int = 1):
    """
    Give an item to another character.
    
    Args:
        item: Name of the item to give
        recipient: Character receiving the item
        quantity: How many to give (default: 1)
    """
    print(f"Gave {quantity}x {item} to {recipient}")
```

The SDK parses the docstring and uses the parameter descriptions in the tool schema sent to the backend.

## Complex Tools

### Tool with Multiple Optional Parameters

```python
@tool
def interact(
    target: str,
    action: str = "talk",
    item: str = None,
    location: str = None
):
    """
    Interact with an object or character.
    
    Args:
        target: What to interact with
        action: Type of interaction (talk, use, examine, etc.)
        item: Optional item to use in the interaction
        location: Optional specific location for the interaction
    """
    if item:
        print(f"{action} with {target} using {item}")
    else:
        print(f"{action} with {target}")
```

### Tool That Returns Values

Tools can return values (though NPCs just see success/failure):

```python
@tool
def search_container(container: str):
    """
    Search a container for items.
    
    Args:
        container: Container to search
    """
    # In a real game, this would check actual container contents
    items = ["gold coins", "health potion"]
    return f"Found: {', '.join(items)}"
```

## Registering Tools

After defining tools, register them with your session:

```python
from llm_npc import NPCClient

client = NPCClient("http://localhost:8080")

with client.session("game-session") as session:
    # Register a list of tool functions
    session.register_tools([
        speak,
        attack,
        cast_spell,
        give_item,
        interact,
        search_container
    ])
    
    # Now create NPCs that can use these tools
    npc = session.create_npc("Wizard", "A powerful wizard")
```

!!! warning "Decorator Required"
    All functions passed to `register_tools()` must be decorated with `@tool`. Otherwise you'll get a `ToolRegistrationError`.

## Tool Best Practices

### 1. Clear Descriptions

Write clear, concise descriptions that tell the NPC **when** to use the tool:

```python
# ✅ Good - clear purpose
@tool(description="Speak dialogue when you want to communicate with someone")
def speak(message: str):
    pass

# ❌ Bad - too vague
@tool(description="A tool for talking")
def speak(message: str):
    pass
```

### 2. Descriptive Parameter Names

Use self-explanatory parameter names:

```python
# ✅ Good
@tool
def move_to(destination: str, movement_speed: str = "walk"):
    pass

# ❌ Bad
@tool
def move_to(dest: str, spd: str = "walk"):
    pass
```

### 3. Document All Parameters

Always include parameter descriptions in docstrings:

```python
@tool
def complex_action(target: str, method: str, intensity: int):
    """
    Perform a complex action.
    
    Args:
        target: Who or what to target
        method: How to perform the action
        intensity: How intense the action should be (1-10)
    """
    pass
```

### 4. Sensible Defaults

Provide sensible defaults for optional parameters:

```python
@tool
def attack(target: str, weapon: str = "fists", power: int = 5):
    """
    Attack a target.
    
    Args:
        target: Who to attack
        weapon: Weapon to use (default: fists for unarmed)
        power: Attack power (default: 5 for normal attack)
    """
    pass
```

### 5. Atomic Actions

Keep tools focused on single actions:

```python
# ✅ Good - single responsibility
@tool
def open_door(door: str):
    """Open a door"""
    pass

@tool
def walk_through_door(door: str):
    """Walk through a door"""
    pass

# ❌ Bad - doing too much
@tool
def open_and_walk_through_door(door: str):
    """Open a door and walk through it"""
    pass
```

## Error Handling

If tool registration fails, you'll get a clear error:

```python
from llm_npc.exceptions import ToolRegistrationError

def not_a_tool():
    """This function isn't decorated"""
    pass

try:
    session.register_tools([not_a_tool])
except ToolRegistrationError as e:
    print(f"Error: {e}")
    # Error: Function not_a_tool is not decorated with @tool
```

## Advanced: Inspecting Tool Metadata

You can inspect the metadata extracted from a tool:

```python
from llm_npc.decorators import get_tool_metadata, is_tool

@tool
def my_tool(param: str):
    """A test tool"""
    pass

# Check if it's a tool
print(is_tool(my_tool))  # True

# Get metadata
metadata = get_tool_metadata(my_tool)
print(metadata)
# {
#     'name': 'my_tool',
#     'description': 'A test tool',
#     'parameters': {
#         'param': {
#             'type': 'string',
#             'description': 'Parameter: param',
#             'required': True
#         }
#     }
# }
```

## Next Steps

- **[Creating NPCs](npcs.md)** - Learn how to create and manage NPCs
- **[Building Context](context.md)** - Provide surroundings and events to NPCs
- **[API Reference: Decorators](../api/decorators.md)** - Full decorator API documentation

