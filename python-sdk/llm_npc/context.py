"""Context builders for constructing NPC surroundings, events, and knowledge graphs."""

from typing import Any, Dict, List, Union
from .models import Surrounding, Event as EventModel


class Surroundings:
    """
    Builder for NPC surroundings.
    
    Provides a simple interface for adding surroundings that can be
    passed directly to NPC.act().
    
    Usage:
        surroundings = Surroundings()
        surroundings.add("Tavern", "A cozy room with wooden tables")
        surroundings.add("Stranger", "A hooded figure in the corner")
    """
    
    def __init__(self):
        self._items: List[Surrounding] = []
    
    def add(self, name: str, description: str) -> "Surroundings":
        """
        Add a surrounding object or entity.
        
        Args:
            name: Name of the object/entity
            description: Description of what it is
        
        Returns:
            Self for method chaining
        """
        self._items.append(Surrounding(name=name, description=description))
        return self
    
    def to_list(self) -> List[Dict[str, str]]:
        """Convert to backend API format (list of dicts)."""
        return [item.to_dict() for item in self._items]
    
    def __iter__(self):
        """Allow iteration over surroundings."""
        return iter(self._items)
    
    def __len__(self):
        """Get number of surroundings."""
        return len(self._items)
    
    def __getitem__(self, index):
        """Get surrounding by index."""
        return self._items[index]


class Event:
    """
    Represents a game event.
    
    Usage:
        event = Event("discovery", "You found a hidden passage")
    """
    
    def __init__(self, event_type: str, event_description: str):
        self.event_type = event_type
        self.event_description = event_description
    
    def to_dict(self) -> Dict[str, str]:
        """Convert to backend API format."""
        return {
            "event_type": self.event_type,
            "event_description": self.event_description
        }


class KnowledgeGraph:
    """
    Builder for knowledge graphs representing NPC memory and relationships.
    
    Usage:
        kg = KnowledgeGraph()
        kg.add_node("player", type="person", name="Hero")
        kg.add_node("quest", type="quest", name="Save the village")
        kg.add_edge("player", "quest", relationship="active")
    """
    
    def __init__(self):
        self._nodes: List[Dict[str, Any]] = []
        self._edges: List[Dict[str, Any]] = []
    
    def add_node(self, node_id: str, **data) -> "KnowledgeGraph":
        """
        Add a node to the knowledge graph.
        
        Args:
            node_id: Unique identifier for this node
            **data: Arbitrary data to attach to the node
        
        Returns:
            Self for method chaining
        """
        self._nodes.append({
            "id": node_id,
            "data": data
        })
        return self
    
    def add_edge(self, source: str, target: str, **data) -> "KnowledgeGraph":
        """
        Add an edge (relationship) between two nodes.
        
        Args:
            source: Source node ID
            target: Target node ID
            **data: Arbitrary data to attach to the edge (e.g., relationship type)
        
        Returns:
            Self for method chaining
        """
        self._edges.append({
            "source": source,
            "target": target,
            "data": data
        })
        return self
    
    def to_dict(self) -> Dict[str, List[Dict[str, Any]]]:
        """Convert to backend API format."""
        return {
            "nodes": self._nodes,
            "edges": self._edges
        }
    
    def __len__(self):
        """Get total number of nodes and edges."""
        return len(self._nodes) + len(self._edges)


class ContextBuilder:
    """
    Advanced builder for constructing complete NPC action context.
    
    Provides a fluent interface for building complex context including
    surroundings, events, and knowledge graphs.
    
    Usage:
        context = ContextBuilder()
        context.add_surrounding("Room", "A dark room")
        context.add_event("noise", "You hear a sound")
        context.set_knowledge_graph(kg)
        
        response = npc.act(context)
    """
    
    def __init__(self):
        self._surroundings = Surroundings()
        self._events: List[Event] = []
        self._knowledge_graph: Union[KnowledgeGraph, None] = None
    
    def add_surrounding(self, name: str, description: str) -> "ContextBuilder":
        """
        Add a surrounding object or entity.
        
        Args:
            name: Name of the object/entity
            description: Description of what it is
        
        Returns:
            Self for method chaining
        """
        self._surroundings.add(name, description)
        return self
    
    def add_event(self, event_type: str, event_description: str) -> "ContextBuilder":
        """
        Add a game event.
        
        Args:
            event_type: Type/category of the event
            event_description: Description of what happened
        
        Returns:
            Self for method chaining
        """
        self._events.append(Event(event_type, event_description))
        return self
    
    def set_knowledge_graph(self, knowledge_graph: KnowledgeGraph) -> "ContextBuilder":
        """
        Set the knowledge graph for this context.
        
        Args:
            knowledge_graph: A KnowledgeGraph instance
        
        Returns:
            Self for method chaining
        """
        self._knowledge_graph = knowledge_graph
        return self
    
    def build(self) -> Dict[str, Any]:
        """
        Build the complete context as a dict for the backend API.
        
        Returns:
            Dict with surroundings, events, and knowledge_graph keys
        """
        result = {
            "surroundings": self._surroundings.to_list(),
        }
        
        if self._events:
            result["events"] = [event.to_dict() for event in self._events]
        
        if self._knowledge_graph:
            result["knowledge_graph"] = self._knowledge_graph.to_dict()
        
        return result
    
    @property
    def surroundings(self) -> Surroundings:
        """Get the surroundings builder."""
        return self._surroundings
    
    @property
    def events(self) -> List[Event]:
        """Get the list of events."""
        return self._events
    
    @property
    def knowledge_graph(self) -> Union[KnowledgeGraph, None]:
        """Get the knowledge graph."""
        return self._knowledge_graph

