# Error Handling

Robust error handling ensures your game remains stable even when the backend or LLM encounters issues.

## Exception Hierarchy

All SDK exceptions inherit from `LLMNPCError`:

```python
LLMNPCError (base)
├── BackendConnectionError    # Cannot connect to backend
├── BackendError              # Backend returned an error
├── NPCNotFoundError          # NPC doesn't exist
├── SessionError              # Session-related errors
└── ToolRegistrationError     # Tool registration failed
```

## Common Error Scenarios

### 1. Backend Not Running

```python
from llm_npc import NPCClient
from llm_npc.exceptions import BackendConnectionError

client = NPCClient("http://localhost:8080")

try:
    if not client.health_check():
        print("Backend is not responding")
except BackendConnectionError:
    print("Cannot connect to backend")
    print("Start it with: ./backend --http")
```

### 2. LLM Timeout

```python
from llm_npc.exceptions import BackendError

try:
    response = npc.act(surroundings, events)
except BackendError as e:
    if "timeout" in str(e).lower():
        print("LLM took too long to respond")
        print("Try: smaller context, faster model, or increase timeout")
```

### 3. Tool Registration Failures

```python
from llm_npc.exceptions import ToolRegistrationError

def not_a_tool():
    """Forgot the @tool decorator"""
    pass

try:
    session.register_tools([not_a_tool])
except ToolRegistrationError as e:
    print(f"Tool registration failed: {e}")
    # Fix: Add @tool decorator
```

### 4. NPC Not Found

```python
from llm_npc.exceptions import NPCNotFoundError

try:
    client.delete_npc("invalid-id")
except NPCNotFoundError:
    print("NPC doesn't exist")
```

## Comprehensive Error Handling

```python
from llm_npc import NPCClient, tool
from llm_npc.exceptions import (
    BackendConnectionError,
    BackendError,
    ToolRegistrationError,
    LLMNPCError
)
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@tool
def speak(message: str):
    """Make NPC speak"""
    pass

def main():
    try:
        # Initialize
        client = NPCClient("http://localhost:8080")
        
        # Health check with timeout handling
        try:
            if not client.health_check():
                logger.error("Backend health check failed")
                return
        except BackendConnectionError as e:
            logger.error(f"Cannot connect to backend: {e}")
            logger.info("Make sure backend is running: ./backend --http")
            return
        
        # Create session
        with client.session("game") as session:
            # Register tools
            try:
                session.register_tools([speak])
            except ToolRegistrationError as e:
                logger.error(f"Tool registration failed: {e}")
                return
            
            # Create NPC
            try:
                npc = session.create_npc(
                    "Guard",
                    "A city guard"
                )
            except BackendError as e:
                logger.error(f"Failed to create NPC: {e}")
                return
            
            # Make NPC act
            try:
                response = npc.act(
                    ["Town square"],
                    ["A stranger approaches"]
                )
                
                if response.success:
                    logger.info(f"NPC response: {response.text}")
                else:
                    logger.warning(f"NPC action failed: {response.error}")
                    
            except BackendError as e:
                logger.error(f"Backend error during action: {e}")
                if hasattr(e, 'status_code'):
                    if e.status_code == 503:
                        logger.error("LLM service unavailable")
                    elif e.status_code == 504:
                        logger.error("LLM request timed out")
                    elif e.status_code == 429:
                        logger.error("Rate limit exceeded")
    
    except LLMNPCError as e:
        logger.error(f"SDK error: {e}")
    except Exception as e:
        logger.error(f"Unexpected error: {e}", exc_info=True)

if __name__ == "__main__":
    main()
```

## Retry Logic

Implement retry logic for transient errors:

```python
import time
from llm_npc.exceptions import BackendError

def npc_act_with_retry(npc, surroundings, events, max_retries=3, delay=1.0):
    """Make NPC act with automatic retries"""
    for attempt in range(max_retries):
        try:
            response = npc.act(surroundings, events)
            
            if response.success:
                return response
            else:
                error = response.error or "Unknown error"
                
                # Check if error is retryable
                if "timeout" in error.lower() or "unavailable" in error.lower():
                    if attempt < max_retries - 1:
                        print(f"Attempt {attempt + 1} failed, retrying...")
                        time.sleep(delay * (attempt + 1))  # Exponential backoff
                        continue
                
                # Non-retryable error
                return response
                
        except BackendError as e:
            if "timeout" in str(e).lower() and attempt < max_retries - 1:
                print(f"Timeout on attempt {attempt + 1}, retrying...")
                time.sleep(delay * (attempt + 1))
                continue
            raise
    
    return None  # All retries failed

# Usage
response = npc_act_with_retry(npc, surroundings, events)
if response:
    print(response.text)
else:
    print("Failed after all retries")
```

## Fallback Behavior

Provide fallback when LLM fails:

```python
def npc_act_with_fallback(npc, surroundings, events, fallback_action=None):
    """Make NPC act with fallback behavior"""
    try:
        response = npc.act(surroundings, events)
        
        if response.success:
            return response
        else:
            print(f"NPC action failed: {response.error}")
            if fallback_action:
                return fallback_action()
    
    except BackendError as e:
        print(f"Backend error: {e}")
        if fallback_action:
            return fallback_action()
    
    return None

# Usage with fallback
def guard_fallback():
    """Default guard behavior when LLM fails"""
    from llm_npc.models import Response, ToolCall, Round
    
    # Create a simple fallback response
    return Response(
        success=True,
        text="The guard stands alert, watching carefully.",
        rounds=[],
        tools_used=[],
        raw_data={}
    )

response = npc_act_with_fallback(guard, surroundings, events, guard_fallback)
```

## Logging

Set up proper logging:

```python
import logging
from llm_npc import NPCClient

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('npc_game.log'),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger(__name__)

# Use in your code
try:
    response = npc.act(surroundings)
    logger.info(f"NPC {npc.name} acted successfully")
except Exception as e:
    logger.error(f"Error with NPC {npc.name}: {e}", exc_info=True)
```

## Validation

Validate responses before using them:

```python
def validate_response(response):
    """Validate NPC response quality"""
    if not response:
        return False, "No response"
    
    if not response.success:
        return False, f"Failed: {response.error}"
    
    # Check for empty response
    if not response.text and not response.tools_used:
        return False, "Empty response (model may be too small)"
    
    # Check tool success
    failed_tools = [t for t in response.tools_used if not t.success]
    if failed_tools:
        return False, f"Tools failed: {[t.name for t in failed_tools]}"
    
    return True, "Valid"

# Usage
response = npc.act(surroundings, events)
is_valid, message = validate_response(response)

if is_valid:
    print(response.text)
else:
    print(f"Invalid response: {message}")
```

## Best Practices

1. **Always catch specific exceptions** before generic ones
2. **Log errors** for debugging
3. **Provide user feedback** for connection issues
4. **Implement retries** for transient errors
5. **Have fallback behavior** for critical NPCs
6. **Validate responses** before using them
7. **Use timeouts** to prevent hanging
8. **Monitor error rates** in production

## See Also

- [Exceptions API Reference](../api/exceptions.md)
- [Best Practices](best-practices.md)
- [Client API](../api/client.md)

