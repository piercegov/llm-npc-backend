#!/bin/bash

# Test script for NPC endpoints
# Requires the backend to be running

SOCKET_PATH="/tmp/llm-npc-backend.sock"

echo "Testing NPC endpoints..."

# Function to make HTTP requests via Unix socket
http_request() {
    local method="$1"
    local path="$2"
    local data="$3"
    
    if [ -n "$data" ]; then
        curl -s --unix-socket "$SOCKET_PATH" \
             -X "$method" \
             -H "Content-Type: application/json" \
             -d "$data" \
             "http://localhost$path"
    else
        curl -s --unix-socket "$SOCKET_PATH" \
             -X "$method" \
             "http://localhost$path"
    fi
    echo ""
}

echo "1. Testing NPC registration..."
REGISTER_RESPONSE=$(http_request "POST" "/npc/register" '{
    "name": "Test Guard",
    "background_story": "A vigilant guard who protects the city gates. Always alert and suspicious of strangers."
}')
echo "Register response: $REGISTER_RESPONSE"

# Extract NPC ID from response
NPC_ID=$(echo "$REGISTER_RESPONSE" | grep -o '"npc_id":"[^"]*"' | cut -d'"' -f4)
echo "Extracted NPC ID: $NPC_ID"

echo -e "\n2. Testing NPC list..."
http_request "GET" "/npc/list"

echo -e "\n3. Testing NPC get..."
http_request "GET" "/npc/$NPC_ID"

echo -e "\n4. Testing NPC act..."
http_request "POST" "/npc/act" "{
    \"npc_id\": \"$NPC_ID\",
    \"surroundings\": [
        {
            \"name\": \"City Gate\",
            \"description\": \"A large wooden gate with iron reinforcements. Several travelers are waiting to enter.\"
        },
        {
            \"name\": \"Suspicious Traveler\",
            \"description\": \"A hooded figure carrying a large pack, avoiding eye contact.\"
        }
    ],
    \"knowledge_graph\": {
        \"nodes\": [],
        \"edges\": []
    },
    \"npc_state\": {},
    \"knowledge_graph_depth\": 0,
    \"events\": [
        {
            \"event_type\": \"new_arrival\",
            \"event_description\": \"A suspicious traveler approached the gate\"
        }
    ]
}"

echo -e "\n5. Testing NPC delete..."
http_request "DELETE" "/npc/$NPC_ID"

echo -e "\n6. Testing NPC list after delete..."
http_request "GET" "/npc/list"

echo -e "\nTest complete!"