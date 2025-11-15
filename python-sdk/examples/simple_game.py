#!/usr/bin/env python3
"""
Simple game example demonstrating the LLM NPC Python SDK.

This example shows how much easier it is to use the SDK compared to
manually constructing API requests.

Compare this to: examples/python/game_client.py
"""

from llm_npc import NPCClient, tool, Surroundings, Event, KnowledgeGraph


# Define game tools with simple decorators
@tool(description="Make the NPC speak dialogue aloud to nearby characters")
def speak(message: str, target: str = None):
    """
    Args:
        message: What the NPC should say
        target: Optional specific character to address
    """
    print(f"[SPEAK] {message}" + (f" (to {target})" if target else ""))


@tool(description="Move the NPC to a different location")
def move_to(location: str):
    """
    Args:
        location: The destination location name
    """
    print(f"[MOVE] Moving to {location}")


@tool(description="Give an item to a character")
def give_item(item: str, recipient: str):
    """
    Args:
        item: The item to give
        recipient: Who to give the item to
    """
    print(f"[GIVE] Giving {item} to {recipient}")


def main():
    print("=== LLM NPC SDK - Simple Game Example ===\n")
    
    # Initialize client
    client = NPCClient("http://localhost:8080")
    
    # Health check
    print("1. Checking backend health...")
    if not client.health_check():
        print("❌ Backend is not running! Start it with: ./backend --http")
        return
    print("✅ Backend is running\n")
    
    # Create session and register tools (much simpler!)
    print("2. Setting up game session...")
    with client.session("simple-game-example") as session:
        # Register tools - SDK automatically converts decorators to backend format
        session.register_tools([speak, move_to, give_item])
        print("✅ Tools registered\n")
        
        # Create NPCs with simple API
        print("3. Creating NPCs...")
        innkeeper = session.create_npc(
            name="Elara the Innkeeper",
            background="A warm innkeeper who runs 'The Gilded Swan' tavern. "
                      "She knows all the local gossip and helps travelers."
        )
        
        guard = session.create_npc(
            name="Captain Marcus",
            background="A veteran city guard captain. Fair but strict, "
                      "and currently tracking a group of thieves."
        )
        print(f"✅ Created: {innkeeper}")
        print(f"✅ Created: {guard}\n")
        
        # Scenario: Suspicious stranger at the inn
        print("4. Running scenario: 'Suspicious Stranger at the Inn'\n")
        print("--- Turn 1: Innkeeper's Perspective ---")
        
        # Build surroundings with simple builder
        surroundings = Surroundings()
        surroundings.add("Tavern Common Room", 
                        "A cozy room with wooden tables, a roaring fireplace, "
                        "and the smell of fresh bread. 5 patrons drinking.")
        surroundings.add("Hooded Stranger",
                        "A figure in a dark cloak sits alone in the corner, "
                        "watching the door intently. Hasn't touched their drink.")
        surroundings.add("Captain Marcus",
                        "The city guard captain just entered through the front door, "
                        "scanning the room.")
        
        # Define events
        events = [
            Event("new_customer", "A hooded stranger entered 10 minutes ago, watching nervously"),
            Event("guard_arrival", "Captain Marcus, the guard captain, just walked in")
        ]
        
        # Execute action - SDK handles all the conversion
        response = innkeeper.act(surroundings, events)
        
        if response.success:
            if response.text:
                print(f"💭 Thought: {response.text}\n")
            
            print(f"Inference rounds: {len(response.rounds)}")
            
            # Check tools used
            if response.tools_used:
                print("🛠️  Tools used:")
                for tool_call in response.tools_used:
                    print(f"  - {tool_call.name}({tool_call.args})")
                    if tool_call.response:
                        print(f"    → {tool_call.response}")
            
            if not response.text and not response.tools_used:
                print("⚠️  Empty response - try a larger model (llama3:8b)")
        else:
            print(f"❌ Error: {response.error}")
        print()
        
        # Turn 2: Guard's perspective
        print("--- Turn 2: Guard's Perspective ---")
        
        # Can also use simple lists (SDK converts automatically)
        guard_surroundings = [
            "The Gilded Swan Tavern - a well-maintained inn",
            "Hooded Figure - suspicious person matching thief description",
            "Elara the innkeeper behind the bar",
            "Several regular customers drinking"
        ]
        
        guard_events = [
            "You spotted someone matching the thieves' description"
        ]
        
        response = guard.act(guard_surroundings, guard_events)
        
        if response.success:
            if response.text:
                print(f"💭 Thought: {response.text}\n")
            
            if response.tools_used:
                print("🛠️  Tools used:")
                for tool_call in response.tools_used:
                    print(f"  - {tool_call.name}({tool_call.args})")
        else:
            print(f"❌ Error: {response.error}")
        print()
        
        # Advanced: Using knowledge graph
        print("5. Using knowledge graph for NPC memory...")
        
        # Build a knowledge graph
        kg = KnowledgeGraph()
        kg.add_node("stranger_01", type="person", name="Hooded Stranger", suspicious=True)
        kg.add_node("guard_marcus", type="person", name="Captain Marcus", role="guard")
        kg.add_node("theft_incident", type="event", description="Series of thefts", date="past week")
        kg.add_edge("guard_marcus", "theft_incident", relationship="investigating")
        kg.add_edge("stranger_01", "theft_incident", relationship="possibly_related")
        
        # New event
        new_events = [
            Event("confrontation", "Captain Marcus approached the hooded stranger's table")
        ]
        
        # NPC with memory context
        response = innkeeper.act(surroundings, new_events, knowledge_graph=kg)
        
        if response.success and response.text:
            print(f"💭 Response with memory: {response.text}")
        print()
    
    print("=== Example Complete ===")
    print("\nNotice how much simpler this is compared to the manual API!")
    print("- @tool decorator instead of manual schema construction")
    print("- Surroundings() builder instead of dict lists")
    print("- Event() helper instead of manual dicts")
    print("- Typed Response object with easy access to tools_used")


if __name__ == "__main__":
    main()

