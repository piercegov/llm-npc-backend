# Creating NPCs

NPCs (Non-Player Characters) are the core of your game's AI. This guide shows you how to create, manage, and interact with NPCs using the SDK.

## Creating an NPC

NPCs are created within a session:

```python
from llm_npc import NPCClient

client = NPCClient("http://localhost:8080")

with client.session("my-game") as session:
    npc = session.create_npc(
        name="Elara the Innkeeper",
        background="A friendly innkeeper who runs The Gilded Swan tavern"
    )
    
    print(npc)  # NPC(id=abc-123, name=Elara the Innkeeper)
```

### Parameters

- **`name`** (required): The NPC's name - used for identification
- **`background`** (required): Background story that shapes the NPC's personality and behavior

## The Background Story

The background story is crucial - it tells the NPC who they are and how to behave.

### Good Background Examples

```python
# ✅ Detailed, defines personality and role
guard = session.create_npc(
    name="Captain Marcus",
    background="A veteran city guard captain with 20 years of experience. "
               "He's fair but strict, takes his duty seriously, and has seen "
               "enough trouble to be cautious around suspicious characters. "
               "Currently tracking a group of thieves."
)

# ✅ Includes relationships and motivations
merchant = session.create_npc(
    name="Tobias the Merchant",
    background="A shrewd merchant who sells rare goods. He's friendly to "
               "regular customers but will haggle ruthlessly with strangers. "
               "He owes money to the thieves' guild and is nervous about it."
)
```

### Background Best Practices

1. **Include personality traits**: Friendly, suspicious, greedy, brave, etc.
2. **Define their role**: What do they do? Why are they here?
3. **Add context**: Current situation, relationships, goals
4. **Be specific**: Vague backgrounds lead to generic behavior

```python
# ❌ Too vague
npc = session.create_npc(
    name="Guard",
    background="A guard"
)

# ✅ Much better
npc = session.create_npc(
    name="Guard Thompson",
    background="A young castle guard on his first week of duty. "
               "Eager to prove himself but nervous about making mistakes. "
               "Takes orders from Captain Marcus."
)
```

## Making NPCs Act

Once created, make NPCs respond to situations with `.act()`:

```python
response = npc.act(
    surroundings=["Town square", "Market stalls", "Fountain"],
    events=["A thief just stole from a vendor"]
)

if response.success:
    print(response.text)
    for tool in response.tools_used:
        print(f"Action: {tool.name}")
```

See [Building Context](context.md) for details on surroundings and events.

## Multiple NPCs

Create as many NPCs as needed:

```python
with client.session("my-game") as session:
    session.register_tools([speak, move_to])
    
    # Create multiple NPCs
    innkeeper = session.create_npc(
        name="Elara",
        background="Friendly innkeeper"
    )
    
    guard = session.create_npc(
        name="Marcus",
        background="Strict guard captain"
    )
    
    thief = session.create_npc(
        name="Shadow",
        background="Stealthy thief"
    )
    
    # Each NPC acts independently
    innkeeper_response = innkeeper.act([...])
    guard_response = guard.act([...])
    thief_response = thief.act([...])
```

## NPC Properties

Each NPC object has these properties:

```python
npc = session.create_npc("Test", "Background")

print(npc.npc_id)      # Backend ID: "abc-123"
print(npc.name)        # "Test"
print(npc.background)  # "Background"
print(npc.session)     # Session object
print(npc.client)      # NPCClient object
```

## Managing NPCs

### List All NPCs

```python
all_npcs = client.list_npcs()
print(f"Total NPCs: {all_npcs['count']}")

for npc_id, info in all_npcs['npcs'].items():
    print(f"- {info['name']} ({npc_id})")
```

### Delete an NPC

```python
# Delete by ID
client.delete_npc(npc.npc_id)

# Or pass the ID directly
client.delete_npc("abc-123")
```

!!! warning "Deletion is Permanent"
    Deleted NPCs cannot be recovered. They're removed from the backend immediately.

## Advanced NPC Patterns

### NPC with Rich Context

Combine NPCs with knowledge graphs for persistent memory:

```python
from llm_npc import KnowledgeGraph

# Create knowledge graph
kg = KnowledgeGraph()
kg.add_node("player", type="person", name="Hero", trust="high")
kg.add_node("quest_01", type="quest", status="active")
kg.add_edge("player", "quest_01", relationship="accepted")

# NPC remembers this context
response = npc.act(
    surroundings=["Castle throne room", "The king"],
    events=["You returned from the quest"],
    knowledge_graph=kg
)
```

See [Knowledge Graphs](knowledge-graphs.md) for more details.

### Contextual NPCs

Same NPC, different contexts:

```python
guard = session.create_npc(
    name="Guard",
    background="City guard on patrol"
)

# Daytime patrol
day_response = guard.act(
    surroundings=["Busy market", "Many shoppers"],
    events=["Everything seems peaceful"]
)

# Nighttime patrol
night_response = guard.act(
    surroundings=["Dark alley", "Suspicious figure"],
    events=["You hear breaking glass"]
)
```

### NPCs with Persistent State

Use scratchpad tools to give NPCs persistent memory across actions:

```python
@tool
def remember(key: str, value: str):
    """
    Store something in memory.
    
    Args:
        key: What to remember
        value: The information
    """
    # Backend handles persistent storage
    pass

@tool
def recall(key: str):
    """
    Recall something from memory.
    
    Args:
        key: What to recall
    """
    # Backend retrieves stored information
    pass

session.register_tools([remember, recall])
npc = session.create_npc("Memory NPC", "An NPC with memory")

# First interaction
response1 = npc.act(["You meet someone", "They say their name is Alice"])
# NPC might use: remember("met_person", "Alice")

# Later interaction
response2 = npc.act(["You see the same person again"])
# NPC might use: recall("met_person") and remember meeting Alice
```

## Error Handling

Handle NPC-related errors:

```python
from llm_npc.exceptions import BackendError, NPCNotFoundError

try:
    npc = session.create_npc("Test", "Background")
except BackendError as e:
    print(f"Failed to create NPC: {e}")
    if e.status_code:
        print(f"HTTP Status: {e.status_code}")

try:
    client.delete_npc("nonexistent-id")
except NPCNotFoundError:
    print("NPC not found")
```

## Session Management

NPCs belong to sessions. When using context managers, cleanup is automatic:

```python
# ✅ Recommended: Context manager
with client.session("game-123") as session:
    npc = session.create_npc("Test", "Background")
    # Use npc...
# Session automatically closed

# ⚠️ Manual management (less common)
session = client.session("game-123")
npc = session.create_npc("Test", "Background")
# Remember to clean up manually if needed
```

## Best Practices

### 1. Descriptive Names

Use descriptive names that help identify NPCs:

```python
# ✅ Good
captain = session.create_npc("Captain Sarah Chen", "...")
merchant = session.create_npc("Old Man Jenkins", "...")

# ❌ Bad
npc1 = session.create_npc("NPC1", "...")
char = session.create_npc("Character", "...")
```

### 2. Detailed Backgrounds

More detail = better NPC behavior:

```python
wizard = session.create_npc(
    name="Aldrin the Wise",
    background="An elderly wizard who served as the king's advisor for 40 years. "
               "He's knowledgeable about ancient magic and history, speaks slowly "
               "and deliberately, and has a habit of quoting old proverbs. "
               "Recently retired and now runs a magic shop in the capital. "
               "He's patient with beginners but won't tolerate disrespect for magic."
)
```

### 3. Consistent Tool Sets

All NPCs in a session share the same tools, so design tools that make sense for all NPCs:

```python
# Register general-purpose tools
session.register_tools([
    speak,           # Everyone can talk
    move_to,         # Everyone can move
    pick_up,         # Everyone can pick things up
    use_item,        # Everyone can use items
])

# Then create NPCs with different personalities
guard = session.create_npc("Guard", "Strict and dutiful")
thief = session.create_npc("Thief", "Sneaky and opportunistic")

# They use the same tools differently based on their backgrounds
```

### 4. Store NPC References

Keep references to NPCs you'll reuse:

```python
class GameNPCs:
    def __init__(self, session):
        self.innkeeper = session.create_npc("Elara", "...")
        self.guard = session.create_npc("Marcus", "...")
        self.merchant = session.create_npc("Tobias", "...")

with client.session("game") as session:
    npcs = GameNPCs(session)
    
    # Easy to reference later
    response = npcs.innkeeper.act([...])
    response = npcs.guard.act([...])
```

## Next Steps

- **[Building Context](context.md)** - Learn how to provide rich context to NPCs
- **[Working with Responses](responses.md)** - Parse and use NPC responses
- **[Knowledge Graphs](knowledge-graphs.md)** - Give NPCs persistent memory
- **[API Reference: Client](../api/client.md)** - Full NPC API documentation

