#!/usr/bin/env python3
"""
Example game client that demonstrates using the LLM NPC Backend.

This is a contrived example showing how a game engine would integrate
with the backend to create intelligent NPCs.
"""

import requests
import json
import time
from typing import Dict, List, Optional

class NPCBackendClient:
    """Client for communicating with the LLM NPC Backend via HTTP."""
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session_id = None
        self.npcs = {}  # Store registered NPC IDs
        
    def health_check(self) -> bool:
        """Check if the backend is running."""
        try:
            response = requests.get(f"{self.base_url}/health")
            return response.text == "pong"
        except Exception as e:
            print(f"Health check failed: {e}")
            return False
    
    def register_tools(self, session_id: str, tools: List[Dict]) -> bool:
        """Register custom game-specific tools for NPCs to use."""
        self.session_id = session_id
        
        payload = {
            "session_id": session_id,
            "tools": tools
        }
        
        response = requests.post(
            f"{self.base_url}/tools/register",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 201:
            result = response.json()
            print(f"✓ Registered {result['tools_count']} tools: {result['tool_names']}")
            return True
        else:
            print(f"✗ Tool registration failed: {response.text}")
            return False
    
    def register_npc(self, name: str, background_story: str) -> Optional[str]:
        """Register a new NPC with the backend."""
        payload = {
            "name": name,
            "background_story": background_story
        }
        
        response = requests.post(
            f"{self.base_url}/npc/register",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 201:
            result = response.json()
            npc_id = result['npc_id']
            self.npcs[name] = npc_id
            print(f"✓ Registered NPC '{name}' with ID: {npc_id}")
            return npc_id
        else:
            print(f"✗ NPC registration failed: {response.text}")
            return None
    
    def npc_act(self, npc_id: str, surroundings: List[Dict], 
                events: List[Dict] = None, knowledge_graph: Dict = None) -> Dict:
        """Execute a tick/action for an NPC."""
        payload = {
            "npc_id": npc_id,
            "surroundings": surroundings,
        }
        
        # Add optional session_id for custom tools
        if self.session_id:
            payload["session_id"] = self.session_id
        
        # Add optional events
        if events:
            payload["events"] = events
        
        # Add optional knowledge graph
        if knowledge_graph:
            payload["knowledge_graph"] = knowledge_graph
        
        response = requests.post(
            f"{self.base_url}/npc/act",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            return response.json()
        else:
            print(f"✗ NPC action failed: {response.status_code} - {response.text}")
            return {"success": False, "error": response.text}
    
    def list_npcs(self) -> Dict:
        """List all registered NPCs."""
        response = requests.get(f"{self.base_url}/npc/list")
        
        if response.status_code == 200:
            return response.json()
        else:
            print(f"✗ Failed to list NPCs: {response.text}")
            return {}
    
    def delete_npc(self, npc_id: str) -> bool:
        """Delete an NPC."""
        response = requests.delete(f"{self.base_url}/npc/{npc_id}")
        
        if response.status_code == 200:
            print(f"✓ Deleted NPC: {npc_id}")
            return True
        else:
            print(f"✗ Failed to delete NPC: {response.text}")
            return False


def main():
    """Example game scenario demonstrating backend usage."""
    print("=== LLM NPC Backend - Example Game Client ===\n")
    
    print("📝 Note: This example works best with larger language models.")
    print("   Small models like qwen3:1.7b may produce empty responses.")
    print("   For better results, use: ollama pull llama3:8b\n")
    
    # Initialize client
    client = NPCBackendClient()
    
    # 1. Health check
    print("1. Checking backend health...")
    if not client.health_check():
        print("Backend is not running! Start it with: ./backend --http")
        return
    print("✓ Backend is running\n")
    
    # 2. Register custom game tools
    print("2. Registering custom game tools...")
    game_tools = [
        {
            "name": "speak",
            "description": "Make the NPC speak dialogue aloud to nearby characters",
            "parameters": {
                "message": {
                    "type": "string",
                    "description": "What the NPC should say",
                    "required": True
                },
                "target": {
                    "type": "string",
                    "description": "Optional: specific character to address",
                    "required": False
                }
            }
        },
        {
            "name": "move_to",
            "description": "Move the NPC to a different location in the game world",
            "parameters": {
                "location": {
                    "type": "string",
                    "description": "The destination location name",
                    "required": True
                }
            }
        },
        {
            "name": "give_item",
            "description": "Give an item from the NPC's inventory to a character",
            "parameters": {
                "item": {
                    "type": "string",
                    "description": "The item to give",
                    "required": True
                },
                "recipient": {
                    "type": "string",
                    "description": "Who to give the item to",
                    "required": True
                }
            }
        }
    ]
    
    if not client.register_tools("example-game-session-001", game_tools):
        return
    print()
    
    # 3. Register NPCs
    print("3. Registering NPCs...")
    
    innkeeper_id = client.register_npc(
        name="Elara the Innkeeper",
        background_story="A warm and welcoming innkeeper who has been running 'The Gilded Swan' tavern for over 20 years. She knows all the local gossip and helps travelers find their way. She's protective of her regulars and suspicious of strangers who cause trouble."
    )
    
    guard_id = client.register_npc(
        name="Captain Marcus",
        background_story="A veteran city guard captain who takes his duty seriously. He's fair but strict, and has seen enough trouble to be cautious around suspicious characters. He's been tracking a group of thieves for weeks."
    )
    print()
    
    # 4. List all NPCs
    print("4. Listing registered NPCs...")
    all_npcs = client.list_npcs()
    print(f"Total NPCs registered: {all_npcs.get('count', 0)}")
    for npc_id, info in all_npcs.get('npcs', {}).items():
        print(f"  - {info['name']} ({npc_id})")
    print()
    
    # 5. Simulate a game scenario
    print("5. Simulating game scenario: 'Suspicious Stranger at the Inn'\n")
    
    # Turn 1: Innkeeper notices a hooded stranger
    print("--- TURN 1: Innkeeper's Perspective ---")
    innkeeper_surroundings = [
        {
            "name": "Tavern Common Room",
            "description": "A cozy room with wooden tables, a roaring fireplace, and the smell of fresh bread. Currently 5 patrons drinking and chatting."
        },
        {
            "name": "Hooded Stranger",
            "description": "A figure in a dark cloak sits alone in the corner, watching the door intently. They haven't touched their drink."
        },
        {
            "name": "Captain Marcus",
            "description": "The city guard captain just entered through the front door, scanning the room."
        }
    ]
    
    innkeeper_events = [
        {
            "event_type": "new_customer",
            "event_description": "A hooded stranger entered 10 minutes ago and has been watching the door nervously"
        },
        {
            "event_type": "guard_arrival",
            "event_description": "Captain Marcus, the city guard captain, just walked in"
        }
    ]
    
    result = client.npc_act(innkeeper_id, innkeeper_surroundings, innkeeper_events)
    
    if result.get('success'):
        response_text = result.get('llm_response', '').strip()
        if response_text:
            print(f"Response: {response_text}\n")
        else:
            print("Response: (empty - NPC may have only used tools or model produced no output)\n")
        
        print(f"Inference rounds: {len(result.get('rounds', []))}")
        
        # Check if NPC used any tools
        for round_num, round_data in enumerate(result.get('rounds', []), 1):
            tools_used = round_data.get('tools_used', [])
            if tools_used:
                print(f"Round {round_num} - Tools used:")
                for tool in tools_used:
                    print(f"  - {tool['tool_name']}: {tool.get('args', {})}")
                    if tool.get('success'):
                        print(f"    → {tool.get('response', 'Success')}")
        
        # Show helpful message if everything is empty
        if not response_text and not any(round_data.get('tools_used') for round_data in result.get('rounds', [])):
            print("⚠️  Note: LLM produced no output. This can happen with small models like qwen3:1.7b.")
            print("    Try using a larger model (llama3:8b, mistral:7b) for better results.")
    else:
        print(f"✗ Error: {result.get('error', 'Unknown error')}")
    print()
    
    # Turn 2: Guard's perspective on the same situation
    print("--- TURN 2: Guard's Perspective ---")
    guard_surroundings = [
        {
            "name": "The Gilded Swan Tavern",
            "description": "A well-maintained inn that's popular with locals. The innkeeper Elara is behind the bar."
        },
        {
            "name": "Hooded Figure",
            "description": "A suspicious person in a dark cloak sitting in the corner. Matches the description of one of the thieves I've been tracking."
        },
        {
            "name": "Other Patrons",
            "description": "Several regular customers drinking and talking, unaware of any danger."
        }
    ]
    
    guard_events = [
        {
            "event_type": "suspect_spotted",
            "event_description": "You spotted someone matching the description of the thieves you've been tracking"
        }
    ]
    
    result = client.npc_act(guard_id, guard_surroundings, guard_events)
    
    if result.get('success'):
        response_text = result.get('llm_response', '').strip()
        if response_text:
            print(f"Response: {response_text}\n")
        else:
            print("Response: (empty - NPC may have only used tools or model produced no output)\n")
        
        print(f"Inference rounds: {len(result.get('rounds', []))}")
        
        # Check if NPC used any tools
        for round_num, round_data in enumerate(result.get('rounds', []), 1):
            tools_used = round_data.get('tools_used', [])
            if tools_used:
                print(f"Round {round_num} - Tools used:")
                for tool in tools_used:
                    print(f"  - {tool['tool_name']}: {tool.get('args', {})}")
                    if tool.get('success'):
                        print(f"    → {tool.get('response', 'Success')}")
        
        # Show helpful message if everything is empty
        if not response_text and not any(round_data.get('tools_used') for round_data in result.get('rounds', [])):
            print("⚠️  Note: LLM produced no output. This can happen with small models like qwen3:1.7b.")
            print("    Try using a larger model (llama3:8b, mistral:7b) for better results.")
    else:
        print(f"✗ Error: {result.get('error', 'Unknown error')}")
    print()
    
    # 6. Demonstrate multi-round thinking with knowledge graph
    print("6. Testing NPC with knowledge graph (memory)...")
    
    # Create a simple knowledge graph representing what the innkeeper knows
    knowledge_graph = {
        "nodes": [
            {"id": "stranger_01", "data": {"type": "person", "name": "Hooded Stranger", "suspicious": True}},
            {"id": "guard_marcus", "data": {"type": "person", "name": "Captain Marcus", "role": "guard"}},
            {"id": "theft_incident", "data": {"type": "event", "description": "Series of thefts in the city", "date": "past week"}}
        ],
        "edges": [
            {"source": "guard_marcus", "target": "theft_incident", "data": {"relationship": "investigating"}},
            {"source": "stranger_01", "target": "theft_incident", "data": {"relationship": "possibly_related"}}
        ]
    }
    
    innkeeper_events_2 = [
        {
            "event_type": "confrontation",
            "event_description": "Captain Marcus approached the hooded stranger's table"
        }
    ]
    
    result = client.npc_act(
        innkeeper_id, 
        innkeeper_surroundings, 
        innkeeper_events_2,
        knowledge_graph
    )
    
    if result.get('success'):
        response_text = result.get('llm_response', '').strip()
        if response_text:
            print(f"Response with KG: {response_text}\n")
        else:
            print("Response with KG: (empty - NPC may have only used tools or model produced no output)\n")
        
        # Check for tools
        for round_num, round_data in enumerate(result.get('rounds', []), 1):
            tools_used = round_data.get('tools_used', [])
            if tools_used:
                print(f"Round {round_num} - Tools used:")
                for tool in tools_used:
                    print(f"  - {tool['tool_name']}: {tool.get('args', {})}")
        
        # Show helpful message if everything is empty
        if not response_text and not any(round_data.get('tools_used') for round_data in result.get('rounds', [])):
            print("⚠️  Note: LLM produced no output. This can happen with small models like qwen3:1.7b.")
            print("    Try using a larger model (llama3:8b, mistral:7b) for better results.")
    else:
        print(f"✗ Error: {result.get('error', 'Unknown error')}")
    
    print("\n=== Example Complete ===")


if __name__ == "__main__":
    main()

