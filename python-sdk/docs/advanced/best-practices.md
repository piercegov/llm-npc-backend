# Best Practices

Production-ready patterns and tips for building robust NPC systems.

## Tool Design

### Keep Tools Atomic

Each tool should do one thing:

```python
# ✅ Good - single responsibility
@tool
def open_door(door: str):
    """Open a specific door"""
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

### Clear Tool Descriptions

Write descriptions that tell NPCs **when** to use the tool:

```python
# ✅ Good
@tool(description="Speak dialogue when you want to communicate with someone nearby")
def speak(message: str):
    pass

# ❌ Bad
@tool(description="A tool for speech")
def speak(message: str):
    pass
```

### Document All Parameters

```python
@tool
def complex_action(target: str, method: str, power: int = 5):
    """
    Perform a complex action on a target.
    
    Args:
        target: Who or what to target
        method: How to perform the action (attack, defend, assist)
        power: Intensity level from 1-10 (default: 5)
    """
    pass
```

## NPC Background Stories

### Be Specific

```python
# ✅ Good - detailed personality
guard = session.create_npc(
    "Captain Sarah Chen",
    "A 15-year veteran of the city guard with a reputation for fairness "
    "but zero tolerance for crime. She's methodical, observant, and "
    "slightly paranoid after a recent assassination attempt. "
    "Currently investigating a theft ring operating in the market district. "
    "She trusts her instincts and isn't easily fooled."
)

# ❌ Bad - too vague
guard = session.create_npc(
    "Guard",
    "A guard who guards things"
)
```

### Include Current Context

```python
merchant = session.create_npc(
    "Tobias the Merchant",
    "A middle-aged merchant who sells rare artifacts. He's usually jovial "
    "and loves to haggle, but he's been nervous lately because he owes money "
    "to the Thieves' Guild and the payment is overdue. He's desperate for a "
    "big sale and might take risks he normally wouldn't."
)
```

## Context Management

### Scope Appropriately

Keep context focused on what matters:

```python
# ✅ Good - focused, relevant details
surroundings = Surroundings()
surroundings.add("Market Square", "Busy marketplace with 20+ vendors")
surroundings.add("Suspicious Man", "Hooded figure watching the jewelry stall")
surroundings.add("Jewelry Stall", "Unguarded stall with expensive items displayed")

# ❌ Bad - too much irrelevant detail
surroundings = Surroundings()
for i in range(50):  # Don't overwhelm the LLM
    surroundings.add(f"Random Detail {i}", "...")
```

### Dynamic Context

Build context from game state:

```python
class ContextManager:
    def __init__(self, game_state):
        self.game_state = game_state
    
    def get_npc_surroundings(self, npc_location):
        """Build surroundings based on NPC location"""
        surroundings = Surroundings()
        
        # Add location
        location = self.game_state.get_location(npc_location)
        surroundings.add(location.name, location.description)
        
        # Add NPCs at same location
        for other_npc in self.game_state.get_npcs_at(npc_location):
            surroundings.add(other_npc.name, other_npc.appearance)
        
        # Add relevant objects
        for obj in self.game_state.get_objects_at(npc_location):
            surroundings.add(obj.name, obj.description)
        
        return surroundings
```

## Response Processing

### Always Validate

```python
def process_npc_action(npc, response):
    """Process NPC action with validation"""
    if not response or not response.success:
        logger.warning(f"NPC {npc.name} action failed")
        return False
    
    # Process tools
    for tool in response.tools_used:
        if not tool.success:
            logger.warning(f"Tool {tool.name} failed for {npc.name}")
            continue
        
        # Handle each tool type
        if tool.name == "speak":
            handle_speech(npc, tool.args)
        elif tool.name == "move_to":
            handle_movement(npc, tool.args)
    
    return True
```

### Cache When Appropriate

```python
class NPCResponseCache:
    def __init__(self, ttl=300):  # 5 minute TTL
        self.cache = {}
        self.ttl = ttl
    
    def get_cache_key(self, npc_id, surroundings, events):
        """Create cache key from inputs"""
        import hashlib
        import json
        
        data = {
            "npc": npc_id,
            "surroundings": str(surroundings),
            "events": str(events)
        }
        return hashlib.md5(json.dumps(data).encode()).hexdigest()
    
    def get(self, npc_id, surroundings, events):
        """Get cached response if valid"""
        key = self.get_cache_key(npc_id, surroundings, events)
        if key in self.cache:
            cached, timestamp = self.cache[key]
            if time.time() - timestamp < self.ttl:
                return cached
        return None
    
    def set(self, npc_id, surroundings, events, response):
        """Cache a response"""
        key = self.get_cache_key(npc_id, surroundings, events)
        self.cache[key] = (response, time.time())
```

## Knowledge Graphs

### Keep Them Focused

```python
def build_npc_knowledge_graph(npc, game_state):
    """Build knowledge graph with only relevant information"""
    kg = KnowledgeGraph()
    
    # Only add entities the NPC knows about
    for entity_id in npc.known_entities:
        entity = game_state.get_entity(entity_id)
        kg.add_node(entity_id, **entity.to_dict())
    
    # Only add relationships the NPC is aware of
    for rel in npc.known_relationships:
        kg.add_edge(rel.source, rel.target, relationship=rel.type)
    
    return kg
```

### Update Incrementally

```python
class NPCKnowledgeManager:
    def __init__(self):
        self.npc_knowledge = {}  # npc_id -> KnowledgeGraph
    
    def update_knowledge(self, npc_id, new_info):
        """Incrementally update NPC knowledge"""
        if npc_id not in self.npc_knowledge:
            self.npc_knowledge[npc_id] = KnowledgeGraph()
        
        kg = self.npc_knowledge[npc_id]
        
        # Add new entities
        for entity in new_info.get("entities", []):
            kg.add_node(entity["id"], **entity["data"])
        
        # Add new relationships
        for rel in new_info.get("relationships", []):
            kg.add_edge(rel["source"], rel["target"], **rel["data"])
    
    def get_knowledge(self, npc_id):
        """Get NPC's current knowledge graph"""
        return self.npc_knowledge.get(npc_id, KnowledgeGraph())
```

## Performance

### Batch Operations

```python
def process_multiple_npcs_parallel(npcs, shared_surroundings, shared_events):
    """Process multiple NPCs in parallel"""
    import concurrent.futures
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        futures = {
            executor.submit(npc.act, shared_surroundings, shared_events): npc
            for npc in npcs
        }
        
        results = {}
        for future in concurrent.futures.as_completed(futures):
            npc = futures[future]
            try:
                response = future.result(timeout=30)
                results[npc.name] = response
            except Exception as e:
                logger.error(f"Error processing {npc.name}: {e}")
                results[npc.name] = None
        
        return results
```

### Reuse Sessions

```python
class GameSessionManager:
    def __init__(self, client):
        self.client = client
        self.session = None
    
    def get_session(self):
        """Get or create session"""
        if not self.session:
            self.session = self.client.session("game-main")
            self.session.register_tools(self.get_all_tools())
        return self.session
    
    def get_all_tools(self):
        """Get all game tools"""
        return [speak, move_to, attack, defend, use_item]
```

## Error Handling

### Graceful Degradation

```python
def npc_act_safe(npc, surroundings, events):
    """NPC action with graceful degradation"""
    try:
        response = npc.act(surroundings, events)
        if response.success:
            return response
    except Exception as e:
        logger.error(f"NPC action failed: {e}")
    
    # Fallback: simple scripted behavior
    return create_fallback_response(npc)

def create_fallback_response(npc):
    """Create a simple scripted response"""
    from llm_npc.models import Response
    
    # Use NPC's background to determine fallback behavior
    if "guard" in npc.background.lower():
        text = "The guard remains watchful and alert."
    elif "merchant" in npc.background.lower():
        text = "The merchant tends to their goods."
    else:
        text = f"{npc.name} continues their activities."
    
    return Response(
        success=True,
        text=text,
        rounds=[],
        tools_used=[],
        raw_data={}
    )
```

## Testing

### Unit Test Tools

```python
from llm_npc.decorators import get_tool_metadata

def test_tool_metadata():
    """Test tool metadata extraction"""
    @tool
    def test_action(param: str, count: int = 5):
        """Test action"""
        pass
    
    metadata = get_tool_metadata(test_action)
    assert metadata["name"] == "test_action"
    assert metadata["description"] == "Test action"
    assert "param" in metadata["parameters"]
    assert metadata["parameters"]["param"]["type"] == "string"
    assert metadata["parameters"]["param"]["required"] == True
    assert metadata["parameters"]["count"]["type"] == "integer"
    assert metadata["parameters"]["count"]["required"] == False
```

### Integration Test

```python
def test_npc_integration():
    """Test complete NPC flow"""
    from llm_npc import NPCClient, tool, Surroundings
    
    @tool
    def speak(message: str):
        pass
    
    client = NPCClient("http://localhost:8080")
    
    # Skip if backend not running
    if not client.health_check():
        pytest.skip("Backend not running")
    
    with client.session("test-session") as session:
        session.register_tools([speak])
        npc = session.create_npc("Test NPC", "A test character")
        
        surroundings = Surroundings()
        surroundings.add("Room", "A test room")
        
        response = npc.act(surroundings)
        
        assert response is not None
        assert response.success or response.error is not None
```

## Monitoring

### Track Metrics

```python
class NPCMetrics:
    def __init__(self):
        self.action_count = 0
        self.success_count = 0
        self.failure_count = 0
        self.total_duration = 0
    
    def record_action(self, response, duration):
        """Record action metrics"""
        self.action_count += 1
        self.total_duration += duration
        
        if response and response.success:
            self.success_count += 1
        else:
            self.failure_count += 1
    
    def get_stats(self):
        """Get statistics"""
        return {
            "total_actions": self.action_count,
            "success_rate": self.success_count / max(self.action_count, 1),
            "avg_duration": self.total_duration / max(self.action_count, 1)
        }

# Usage
metrics = NPCMetrics()

start = time.time()
response = npc.act(surroundings, events)
duration = time.time() - start

metrics.record_action(response, duration)
```

## See Also

- [Error Handling](error-handling.md)
- [API Reference](../api/client.md)
- [Examples](../examples/index.md)

