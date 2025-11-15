# Exceptions API Reference

The exceptions module provides custom exception classes for error handling.

## Exception Hierarchy

```
LLMNPCError (base)
├── BackendConnectionError
├── BackendError
├── NPCNotFoundError
├── SessionError
└── ToolRegistrationError
```

## LLMNPCError

::: llm_npc.exceptions.LLMNPCError
    options:
      show_root_heading: true
      show_source: false

## BackendConnectionError

::: llm_npc.exceptions.BackendConnectionError
    options:
      show_root_heading: true
      show_source: false

## BackendError

::: llm_npc.exceptions.BackendError
    options:
      show_root_heading: true
      show_source: false

## NPCNotFoundError

::: llm_npc.exceptions.NPCNotFoundError
    options:
      show_root_heading: true
      show_source: false

## SessionError

::: llm_npc.exceptions.SessionError
    options:
      show_root_heading: true
      show_source: false

## ToolRegistrationError

::: llm_npc.exceptions.ToolRegistrationError
    options:
      show_root_heading: true
      show_source: false

## Usage Examples

### Handling Connection Errors

```python
from llm_npc import NPCClient
from llm_npc.exceptions import BackendConnectionError

client = NPCClient("http://localhost:8080")

try:
    if not client.health_check():
        print("Backend is not running")
except BackendConnectionError as e:
    print(f"Cannot connect to backend: {e}")
    print("Make sure the backend server is running")
```

### Handling Backend Errors

```python
from llm_npc.exceptions import BackendError

try:
    response = npc.act(surroundings, events)
except BackendError as e:
    print(f"Backend error: {e}")
    if e.status_code:
        print(f"HTTP Status Code: {e.status_code}")
```

### Handling Tool Registration Errors

```python
from llm_npc.exceptions import ToolRegistrationError

def not_decorated():
    """This function is missing the @tool decorator"""
    pass

try:
    session.register_tools([not_decorated])
except ToolRegistrationError as e:
    print(f"Tool registration failed: {e}")
    # "Function not_decorated is not decorated with @tool"
```

### Handling NPC Errors

```python
from llm_npc.exceptions import NPCNotFoundError

try:
    client.delete_npc("nonexistent-id")
except NPCNotFoundError as e:
    print(f"NPC not found: {e}")
```

### Catching All SDK Errors

```python
from llm_npc.exceptions import LLMNPCError

try:
    # Any SDK operation
    response = npc.act(surroundings)
except LLMNPCError as e:
    # Catches all SDK exceptions
    print(f"SDK error: {e}")
except Exception as e:
    # Catches other exceptions (network, etc.)
    print(f"Unexpected error: {e}")
```

### Complete Error Handling Example

```python
from llm_npc import NPCClient, tool
from llm_npc.exceptions import (
    BackendConnectionError,
    BackendError,
    ToolRegistrationError,
    SessionError,
    LLMNPCError
)

@tool
def speak(message: str):
    """Make NPC speak"""
    pass

try:
    # Connect to backend
    client = NPCClient("http://localhost:8080")
    
    if not client.health_check():
        raise BackendConnectionError("Backend not responding")
    
    # Create session and NPC
    with client.session("game") as session:
        try:
            session.register_tools([speak])
        except ToolRegistrationError as e:
            print(f"Tool registration failed: {e}")
            raise
        
        npc = session.create_npc("Guard", "A city guard")
        
        # Make NPC act
        try:
            response = npc.act(["Town square"], ["A thief runs past"])
            if not response.success:
                print(f"NPC action failed: {response.error}")
        except BackendError as e:
            print(f"Backend error during action: {e}")
            if e.status_code == 503:
                print("LLM service unavailable")
            elif e.status_code == 504:
                print("LLM timeout")

except BackendConnectionError as e:
    print(f"Cannot connect: {e}")
except SessionError as e:
    print(f"Session error: {e}")
except LLMNPCError as e:
    print(f"SDK error: {e}")
except Exception as e:
    print(f"Unexpected error: {e}")
```

## See Also

- [Error Handling Guide](../advanced/error-handling.md)
- [Client API](client.md)

