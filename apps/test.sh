#!/bin/bash

BASE_URL="http://localhost:8090/api/v1"
HEADER="Content-Type: application/json"

echo "=== 1. Получаем список устройств ==="
DEVICES_JSON=$(curl -s "$BASE_URL/devices")
echo "$DEVICES_JSON" | jq || echo "$DEVICES_JSON"

echo
echo "=== 2. Удаляем все существующие устройства ==="

IDS=$(echo "$DEVICES_JSON" | jq -r '.[].id' 2>/dev/null)

if [ -z "$IDS" ]; then
  echo "Нет устройств для удаления"
else
  for ID in $IDS; do
    echo "Deleting device $ID"
    curl -s -X DELETE "$BASE_URL/devices/$ID" | jq
  done
fi

echo
echo "=== 3. Проверяем, что список пуст ==="
curl -s "$BASE_URL/devices" | jq

echo
echo "=== 4. Создаём устройства ==="

TEMP=$(curl -s -X POST "$BASE_URL/devices" \
  -H "$HEADER" \
  -d '{
    "name": "Temperature Sensor",
    "type": "temperature",
    "location": "kitchen",
    "unit": "C"
  }')
echo "$TEMP" | jq

LIGHT=$(curl -s -X POST "$BASE_URL/devices" \
  -H "$HEADER" \
  -d '{
    "name": "Living Room Light",
    "type": "light",
    "location": "living_room",
    "unit": "boolean"
  }')
echo "$LIGHT" | jq

GATE=$(curl -s -X POST "$BASE_URL/devices" \
  -H "$HEADER" \
  -d '{
    "name": "Main Gate",
    "type": "gate",
    "location": "yard",
    "unit": "boolean"
  }')
echo "$GATE" | jq

echo
echo "=== 5. Получаем список устройств ==="
curl -s "$BASE_URL/devices" | jq

echo
echo "=== 6. Обновляем значения (commands) ==="

curl -s -X POST "$BASE_URL/devices/1/command" \
  -H "$HEADER" \
  -d '{"value": 22.5, "status": "active"}' | jq

curl -s -X POST "$BASE_URL/devices/2/command" \
  -H "$HEADER" \
  -d '{"value": 1, "status": "on"}' | jq

curl -s -X POST "$BASE_URL/devices/3/command" \
  -H "$HEADER" \
  -d '{"value": 0, "status": "closed"}' | jq

echo
echo "=== 7. Проверяем устройства по ID ==="

for ID in 1 2 3; do
  echo "--- Device $ID ---"
  curl -s "$BASE_URL/devices/$ID" | jq
done

echo
echo "✅ DONE: Device Management Service REST test finished"
