"""
LLM NPC Python SDK

A simple and powerful SDK for building intelligent NPCs using the LLM NPC Backend.

Basic usage:
    from llm_npc import NPCClient, tool

    @tool
    def speak(message: str):
        '''Make the NPC speak'''
        pass

    client = NPCClient("http://localhost:8080")
    
    with client.session("my-game") as session:
        session.register_tools([speak])
        npc = session.create_npc("Gandalf", "A wise wizard")
        response = npc.act(surroundings=["Dark cave", "Dragon"])
"""

__version__ = "0.1.0"

from .client import NPCClient, Session, NPC
from .decorators import tool
from .context import Surroundings, Event, ContextBuilder, KnowledgeGraph
from .models import Response, ToolCall, Round, Surrounding, Event as EventModel
from .exceptions import (
    LLMNPCError,
    BackendConnectionError,
    BackendError,
    NPCNotFoundError,
    SessionError,
    ToolRegistrationError
)

__all__ = [
    # Client classes
    "NPCClient",
    "Session",
    "NPC",
    
    # Decorators
    "tool",
    
    # Context builders
    "Surroundings",
    "Event",
    "ContextBuilder",
    "KnowledgeGraph",
    
    # Models
    "Response",
    "ToolCall",
    "Round",
    "Surrounding",
    "EventModel",
    
    # Exceptions
    "LLMNPCError",
    "BackendConnectionError",
    "BackendError",
    "NPCNotFoundError",
    "SessionError",
    "ToolRegistrationError",
]

