"""Typed models for the LLM NPC SDK."""

from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional


@dataclass
class Surrounding:
    """Represents an object or entity in the NPC's surroundings."""
    
    name: str
    description: str
    
    def to_dict(self) -> Dict[str, str]:
        """Convert to backend API format."""
        return {
            "name": self.name,
            "description": self.description
        }


@dataclass
class Event:
    """Represents an event that has occurred in the game."""
    
    event_type: str
    event_description: str
    
    def to_dict(self) -> Dict[str, str]:
        """Convert to backend API format."""
        return {
            "event_type": self.event_type,
            "event_description": self.event_description
        }


@dataclass
class ToolCall:
    """Represents a tool call made by the NPC."""
    
    tool_name: str
    args: Dict[str, Any]
    success: bool
    response: Optional[str] = None
    result: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ToolCall":
        """Create from backend response."""
        return cls(
            tool_name=data.get("tool_name", ""),
            args=data.get("args", {}),
            success=data.get("success", False),
            response=data.get("response"),
            result=data.get("result")
        )
    
    @property
    def name(self) -> str:
        """Alias for tool_name for convenience."""
        return self.tool_name


@dataclass
class Round:
    """Represents a single round of NPC inference."""
    
    tools_used: List[ToolCall] = field(default_factory=list)
    raw_data: Dict[str, Any] = field(default_factory=dict)
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Round":
        """Create from backend response."""
        tools_used = [
            ToolCall.from_dict(tool) 
            for tool in data.get("tools_used", [])
        ]
        return cls(
            tools_used=tools_used,
            raw_data=data
        )


@dataclass
class Response:
    """Represents the response from an NPC action."""
    
    success: bool
    text: str
    rounds: List[Round] = field(default_factory=list)
    tools_used: List[ToolCall] = field(default_factory=list)
    raw_data: Dict[str, Any] = field(default_factory=dict)
    error: Optional[str] = None
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Response":
        """Create from backend response."""
        success = data.get("success", False)
        
        # Parse rounds
        rounds = [
            Round.from_dict(round_data) 
            for round_data in data.get("rounds", [])
        ]
        
        # Flatten all tools used across all rounds for convenience
        all_tools = []
        for round_obj in rounds:
            all_tools.extend(round_obj.tools_used)
        
        # Get the LLM response text
        text = data.get("llm_response", "").strip()
        
        return cls(
            success=success,
            text=text,
            rounds=rounds,
            tools_used=all_tools,
            raw_data=data,
            error=data.get("error")
        )
    
    @property
    def llm_response(self) -> str:
        """Alias for text for compatibility."""
        return self.text

