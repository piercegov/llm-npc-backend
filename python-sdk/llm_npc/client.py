"""Main client for communicating with the LLM NPC Backend."""

from typing import Any, Callable, Dict, List, Optional, Union
import requests

from .models import Response, Surrounding, Event as EventModel
from .context import Surroundings, Event, ContextBuilder, KnowledgeGraph
from .decorators import is_tool, tool_to_backend_format
from .exceptions import (
    BackendConnectionError,
    BackendError,
    NPCNotFoundError,
    SessionError,
    ToolRegistrationError
)


class NPCClient:
    """
    Main client for interacting with the LLM NPC Backend.
    
    Usage:
        client = NPCClient("http://localhost:8080")
        
        with client.session("my-game-session") as session:
            npc = session.create_npc("Gandalf", "A wise wizard")
            response = npc.act(surroundings=["Dark cave", "Dragon"])
    """
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        """
        Initialize the NPC client.
        
        Args:
            base_url: Base URL of the backend server
        """
        self.base_url = base_url.rstrip('/')
        self._session = requests.Session()
    
    def health_check(self) -> bool:
        """
        Check if the backend is running.
        
        Returns:
            True if backend is healthy, False otherwise
        """
        try:
            response = self._session.get(f"{self.base_url}/health", timeout=5)
            return response.text == "pong"
        except Exception:
            return False
    
    def session(self, session_id: str) -> "Session":
        """
        Create a new session for managing tools and NPCs.
        
        Args:
            session_id: Unique identifier for this game session
        
        Returns:
            A Session instance
        """
        return Session(self, session_id)
    
    def list_npcs(self) -> Dict[str, Any]:
        """
        List all registered NPCs.
        
        Returns:
            Dict with 'count' and 'npcs' keys
        """
        try:
            response = self._session.get(f"{self.base_url}/npc/list")
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            raise BackendConnectionError(f"Failed to list NPCs: {e}")
    
    def delete_npc(self, npc_id: str) -> bool:
        """
        Delete an NPC.
        
        Args:
            npc_id: ID of the NPC to delete
        
        Returns:
            True if successful
        """
        try:
            response = self._session.delete(f"{self.base_url}/npc/{npc_id}")
            response.raise_for_status()
            return True
        except requests.RequestException as e:
            raise BackendError(f"Failed to delete NPC: {e}")


class Session:
    """
    Manages a game session with tools and NPCs.
    
    Can be used as a context manager:
        with client.session("my-session") as session:
            session.register_tools([...])
            npc = session.create_npc(...)
    """
    
    def __init__(self, client: NPCClient, session_id: str):
        """
        Initialize a session.
        
        Args:
            client: The NPCClient instance
            session_id: Unique identifier for this session
        """
        self.client = client
        self.session_id = session_id
        self._tools_registered = False
    
    def __enter__(self) -> "Session":
        """Enter context manager."""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Exit context manager."""
        pass
    
    def register_tools(self, tools: List[Callable]) -> "Session":
        """
        Register game-specific tools for NPCs to use.
        
        Args:
            tools: List of functions decorated with @tool
        
        Returns:
            Self for method chaining
        
        Raises:
            ToolRegistrationError: If registration fails
        """
        # Convert tools to backend format
        tool_specs = []
        for tool_func in tools:
            if not is_tool(tool_func):
                raise ToolRegistrationError(
                    f"Function {tool_func.__name__} is not decorated with @tool"
                )
            tool_specs.append(tool_to_backend_format(tool_func))
        
        # Send to backend
        payload = {
            "session_id": self.session_id,
            "tools": tool_specs
        }
        
        try:
            response = self.client._session.post(
                f"{self.client.base_url}/tools/register",
                json=payload,
                headers={"Content-Type": "application/json"}
            )
            response.raise_for_status()
            self._tools_registered = True
            return self
        except requests.RequestException as e:
            raise ToolRegistrationError(f"Failed to register tools: {e}")
    
    def create_npc(self, name: str, background: str) -> "NPC":
        """
        Create a new NPC.
        
        Args:
            name: Name of the NPC
            background: Background story/description
        
        Returns:
            An NPC instance
        
        Raises:
            BackendError: If NPC creation fails
        """
        payload = {
            "name": name,
            "background_story": background
        }
        
        try:
            response = self.client._session.post(
                f"{self.client.base_url}/npc/register",
                json=payload,
                headers={"Content-Type": "application/json"}
            )
            response.raise_for_status()
            result = response.json()
            npc_id = result.get('npc_id')
            
            if not npc_id:
                raise BackendError("Backend did not return npc_id")
            
            return NPC(
                client=self.client,
                session=self,
                npc_id=npc_id,
                name=name,
                background=background
            )
        except requests.RequestException as e:
            raise BackendError(f"Failed to create NPC: {e}")


class NPC:
    """
    Represents an NPC that can perform actions.
    
    Usage:
        response = npc.act(
            surroundings=["Forest", "Sword on ground"],
            events=["You found a weapon"]
        )
    """
    
    def __init__(
        self,
        client: NPCClient,
        session: Session,
        npc_id: str,
        name: str,
        background: str
    ):
        """
        Initialize an NPC instance.
        
        Args:
            client: The NPCClient instance
            session: The Session this NPC belongs to
            npc_id: Backend NPC ID
            name: NPC name
            background: NPC background story
        """
        self.client = client
        self.session = session
        self.npc_id = npc_id
        self.name = name
        self.background = background
    
    def act(
        self,
        surroundings: Union[
            List[str],
            List[Dict[str, str]],
            List[Surrounding],
            Surroundings,
            ContextBuilder
        ],
        events: Optional[Union[List[str], List[Dict[str, str]], List[Event]]] = None,
        knowledge_graph: Optional[Union[KnowledgeGraph, Dict]] = None
    ) -> Response:
        """
        Execute an action/tick for this NPC.
        
        Args:
            surroundings: What the NPC can see/interact with. Can be:
                - Simple list of strings (converted to surroundings)
                - List of dicts with "name" and "description"
                - List of Surrounding objects
                - Surroundings builder
                - ContextBuilder (other args ignored if used)
            events: Recent events that occurred. Can be:
                - Simple list of strings
                - List of dicts with "event_type" and "event_description"
                - List of Event objects
            knowledge_graph: Optional knowledge graph for NPC memory. Can be:
                - KnowledgeGraph object
                - Dict with "nodes" and "edges"
        
        Returns:
            Response object with NPC's action result
        
        Raises:
            BackendError: If the action fails
        """
        # Build the payload
        payload = {"npc_id": self.npc_id}
        
        # Add session_id if tools were registered
        if self.session._tools_registered:
            payload["session_id"] = self.session.session_id
        
        # Handle different input types for surroundings
        if isinstance(surroundings, ContextBuilder):
            # ContextBuilder provides everything
            context = surroundings.build()
            payload["surroundings"] = context["surroundings"]
            if "events" in context:
                payload["events"] = context["events"]
            if "knowledge_graph" in context:
                payload["knowledge_graph"] = context["knowledge_graph"]
        else:
            # Convert surroundings to proper format
            payload["surroundings"] = self._convert_surroundings(surroundings)
            
            # Convert events if provided
            if events is not None:
                payload["events"] = self._convert_events(events)
            
            # Convert knowledge graph if provided
            if knowledge_graph is not None:
                payload["knowledge_graph"] = self._convert_knowledge_graph(knowledge_graph)
        
        # Make the request
        try:
            response = self.client._session.post(
                f"{self.client.base_url}/npc/act",
                json=payload,
                headers={"Content-Type": "application/json"}
            )
            response.raise_for_status()
            result = response.json()
            
            return Response.from_dict(result)
        except requests.RequestException as e:
            raise BackendError(f"Failed to execute NPC action: {e}")
    
    def _convert_surroundings(
        self,
        surroundings: Union[
            List[str],
            List[Dict[str, str]],
            List[Surrounding],
            Surroundings
        ]
    ) -> List[Dict[str, str]]:
        """Convert surroundings to backend format."""
        if isinstance(surroundings, Surroundings):
            return surroundings.to_list()
        
        result = []
        for item in surroundings:
            if isinstance(item, str):
                # Simple string, use as both name and description
                result.append({"name": item, "description": item})
            elif isinstance(item, Surrounding):
                result.append(item.to_dict())
            elif isinstance(item, dict):
                result.append(item)
            else:
                raise ValueError(f"Invalid surrounding type: {type(item)}")
        
        return result
    
    def _convert_events(
        self,
        events: Union[List[str], List[Dict[str, str]], List[Event]]
    ) -> List[Dict[str, str]]:
        """Convert events to backend format."""
        result = []
        for item in events:
            if isinstance(item, str):
                # Simple string, use as both type and description
                result.append({
                    "event_type": "event",
                    "event_description": item
                })
            elif isinstance(item, Event):
                result.append(item.to_dict())
            elif isinstance(item, EventModel):
                result.append(item.to_dict())
            elif isinstance(item, dict):
                result.append(item)
            else:
                raise ValueError(f"Invalid event type: {type(item)}")
        
        return result
    
    def _convert_knowledge_graph(
        self,
        knowledge_graph: Union[KnowledgeGraph, Dict]
    ) -> Dict[str, Any]:
        """Convert knowledge graph to backend format."""
        if isinstance(knowledge_graph, KnowledgeGraph):
            return knowledge_graph.to_dict()
        elif isinstance(knowledge_graph, dict):
            return knowledge_graph
        else:
            raise ValueError(f"Invalid knowledge graph type: {type(knowledge_graph)}")
    
    def __repr__(self) -> str:
        return f"NPC(id={self.npc_id}, name={self.name})"

