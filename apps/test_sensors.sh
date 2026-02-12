#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=== Testing Sensor API ==="

# 1. Create a new sensor
echo "Creating a new sensor..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/sensors" \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Living Room Temp",
        "type": "temperature",
        "location": "living_room"
      }')
echo "Response: $CREATE_RESPONSE"

SENSOR_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
if [ "$SENSOR_ID" = "null" ] || [ -z "$SENSOR_ID" ]; then
  echo "Failed to create sensor. Exiting."
  exit 1
fi
echo "Created sensor with ID: $SENSOR_ID"

# 2. Get all sensors
echo "Fetching all sensors..."
curl -s "$BASE_URL/sensors" | jq

# 3. Get sensor by ID
echo "Fetching sensor by ID..."
curl -s "$BASE_URL/sensors/$SENSOR_ID" | jq

# 4. Update sensor
echo "Updating sensor..."
curl -s -X PUT "$BASE_URL/sensors/$SENSOR_ID" \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Living Room Temperature",
        "location": "living_room_updated"
      }' | jq

# 5. Update sensor value
echo "Updating sensor value..."
curl -s -X PATCH "$BASE_URL/sensors/$SENSOR_ID/value" \
  -H "Content-Type: application/json" \
  -d '{
        "value": 23.5,
        "status": "ok"
      }' | jq

# 6. Get temperature by location
echo "Fetching temperature by location..."
curl -s "$BASE_URL/sensors/temperature/living_room_updated" | jq

# 7. Delete sensor
echo "Deleting sensor..."
curl -s -X DELETE "$BASE_URL/sensors/$SENSOR_ID" | jq

echo "=== All tests completed ==="
