# Decorators API Reference

The decorators module provides the `@tool` decorator for defining game actions.

## tool

::: llm_npc.decorators.tool
    options:
      show_root_heading: true
      show_source: false

## Helper Functions

::: llm_npc.decorators.is_tool
    options:
      show_root_heading: true
      show_source: false

::: llm_npc.decorators.get_tool_metadata
    options:
      show_root_heading: true
      show_source: false

::: llm_npc.decorators.tool_to_backend_format
    options:
      show_root_heading: true
      show_source: false

## Usage Examples

### Basic Tool

```python
from llm_npc import tool

@tool
def speak(message: str):
    """Make the NPC speak dialogue"""
    print(f"NPC says: {message}")
```

### Tool with Custom Description

```python
@tool(description="Attack a target with a weapon")
def attack(target: str, weapon: str = "fists"):
    """
    Args:
        target: Enemy to attack
        weapon: Weapon to use
    """
    print(f"Attacked {target} with {weapon}")
```

### Tool with Multiple Parameters

```python
@tool
def cast_spell(
    spell_name: str,
    target: str = None,
    power: int = 5,
    verbal: bool = True
):
    """
    Cast a magical spell.
    
    Args:
        spell_name: Name of the spell
        target: Optional target
        power: Spell power level
        verbal: Whether to speak the incantation
    """
    if target:
        print(f"Cast {spell_name} on {target} (power: {power})")
    else:
        print(f"Cast {spell_name} (power: {power})")
```

### Inspecting Tool Metadata

```python
from llm_npc.decorators import is_tool, get_tool_metadata

@tool
def my_action(param: str):
    """An example action"""
    pass

# Check if it's a tool
print(is_tool(my_action))  # True

# Get metadata
metadata = get_tool_metadata(my_action)
print(metadata["name"])  # "my_action"
print(metadata["description"])  # "An example action"
print(metadata["parameters"])  # {'param': {...}}
```

## Parameter Type Mapping

| Python Type | Backend Type |
|-------------|--------------|
| `str` | `"string"` |
| `int` | `"integer"` |
| `float` | `"float"` |
| `bool` | `"boolean"` |
| `Optional[T]` | Same as `T` (but not required) |

## See Also

- [Defining Tools Guide](../user-guide/tools.md)
- [Client API](client.md)

