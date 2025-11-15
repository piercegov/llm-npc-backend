"""Decorators for defining game tools."""

import inspect
from typing import Any, Callable, Dict, Optional, get_type_hints
from functools import wraps


def tool(
    func: Optional[Callable] = None,
    *,
    description: Optional[str] = None
) -> Callable:
    """
    Decorator to mark a function as a tool that NPCs can use.
    
    Automatically extracts:
    - Tool name from function name
    - Description from docstring or explicit description parameter
    - Parameters from function signature with type hints
    
    Args:
        func: The function to decorate (when used without arguments)
        description: Optional explicit description (overrides docstring)
    
    Usage:
        @tool
        def speak(message: str, target: str = None):
            '''Make the NPC speak dialogue'''
            pass
        
        @tool(description="Custom description")
        def attack(target: str, damage: int = 10):
            pass
    """
    def decorator(f: Callable) -> Callable:
        # Extract metadata
        tool_name = f.__name__
        tool_description = description or _extract_description(f)
        parameters = _extract_parameters(f)
        
        # Store metadata on the function
        f._tool_metadata = {
            "name": tool_name,
            "description": tool_description,
            "parameters": parameters
        }
        
        # Mark it as a tool
        f._is_tool = True
        
        @wraps(f)
        def wrapper(*args, **kwargs):
            return f(*args, **kwargs)
        
        # Copy metadata to wrapper
        wrapper._tool_metadata = f._tool_metadata
        wrapper._is_tool = True
        
        return wrapper
    
    # Handle both @tool and @tool() syntax
    if func is None:
        # Called with arguments: @tool(description="...")
        return decorator
    else:
        # Called without arguments: @tool
        return decorator(func)


def _extract_description(func: Callable) -> str:
    """Extract description from function docstring."""
    docstring = inspect.getdoc(func)
    if not docstring:
        return f"Tool: {func.__name__}"
    
    # Get first line of docstring as description
    lines = docstring.strip().split('\n')
    return lines[0].strip()


def _extract_parameters(func: Callable) -> Dict[str, Dict[str, Any]]:
    """
    Extract parameter information from function signature.
    
    Returns a dict matching the backend tool schema format:
    {
        "param_name": {
            "type": "string|integer|boolean|float",
            "description": "...",
            "required": true|false
        }
    }
    """
    sig = inspect.signature(func)
    type_hints = get_type_hints(func) if hasattr(inspect, 'get_type_hints') else {}
    
    # Parse docstring for parameter descriptions
    param_descriptions = _parse_param_descriptions(func)
    
    parameters = {}
    
    for param_name, param in sig.parameters.items():
        if param_name == 'self':
            continue
        
        # Determine if required (no default value)
        required = param.default == inspect.Parameter.empty
        
        # Get type from type hints
        param_type = "string"  # default
        if param_name in type_hints:
            py_type = type_hints[param_name]
            param_type = _python_type_to_schema_type(py_type)
        
        # Get description from docstring
        param_desc = param_descriptions.get(
            param_name,
            f"Parameter: {param_name}"
        )
        
        parameters[param_name] = {
            "type": param_type,
            "description": param_desc,
            "required": required
        }
    
    return parameters


def _parse_param_descriptions(func: Callable) -> Dict[str, str]:
    """Parse parameter descriptions from docstring."""
    docstring = inspect.getdoc(func)
    if not docstring:
        return {}
    
    descriptions = {}
    lines = docstring.split('\n')
    
    # Look for Args: section
    in_args_section = False
    for line in lines:
        stripped = line.strip()
        
        if stripped.lower().startswith('args:') or stripped.lower().startswith('parameters:'):
            in_args_section = True
            continue
        
        if in_args_section:
            # Stop at next section
            if stripped and not stripped[0].isspace() and ':' in stripped and stripped[0].isupper():
                break
            
            # Parse parameter line: "param_name: description" or "param_name (type): description"
            if ':' in stripped:
                parts = stripped.split(':', 1)
                param_part = parts[0].strip()
                desc_part = parts[1].strip() if len(parts) > 1 else ""
                
                # Remove type hints in parentheses
                if '(' in param_part:
                    param_part = param_part.split('(')[0].strip()
                
                if param_part:
                    descriptions[param_part] = desc_part
    
    return descriptions


def _python_type_to_schema_type(py_type: Any) -> str:
    """Convert Python type hint to backend schema type string."""
    # Handle Optional types
    origin = getattr(py_type, '__origin__', None)
    if origin is not None:
        # Handle Optional[T] which is Union[T, None]
        args = getattr(py_type, '__args__', ())
        if origin is type(None) or (hasattr(origin, '__name__') and 'Union' in str(origin)):
            # Get the non-None type
            non_none_types = [t for t in args if t is not type(None)]
            if non_none_types:
                py_type = non_none_types[0]
    
    # Map Python types to schema types
    if py_type == str or py_type == 'str':
        return "string"
    elif py_type == int or py_type == 'int':
        return "integer"
    elif py_type == float or py_type == 'float':
        return "float"
    elif py_type == bool or py_type == 'bool':
        return "boolean"
    else:
        return "string"  # default


def is_tool(func: Callable) -> bool:
    """Check if a function is decorated with @tool."""
    return getattr(func, '_is_tool', False)


def get_tool_metadata(func: Callable) -> Optional[Dict[str, Any]]:
    """Get the tool metadata from a decorated function."""
    return getattr(func, '_tool_metadata', None)


def tool_to_backend_format(func: Callable) -> Dict[str, Any]:
    """
    Convert a decorated tool function to backend API format.
    
    Returns:
        {
            "name": "tool_name",
            "description": "Tool description",
            "parameters": {
                "param1": {"type": "string", "description": "...", "required": true},
                ...
            }
        }
    """
    if not is_tool(func):
        raise ValueError(f"Function {func.__name__} is not decorated with @tool")
    
    return get_tool_metadata(func)

