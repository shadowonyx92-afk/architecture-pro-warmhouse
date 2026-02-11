#!/bin/bash

set -e

API_URL="http://localhost:8090/api/v1"
TELEMETRY_URL="http://localhost:8100"

# echo "Очищаем таблицу devices..."
# psql "$DB_URL" -c "TRUNCATE TABLE devices RESTART IDENTITY CASCADE;"
# echo "База очищена."

echo "Добавляем тестовые устройства..."

# Temperature Device
curl -s -X POST "$API_URL/devices" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Temperature Sensor",
        "type": "temperature",
        "location": "Living Room",
        "value": "33.5",
        "unit": "C",
        "status": "active"
    }' | jq
echo

# Light Device
curl -s -X POST "$API_URL/devices" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Light",
        "type": "light",
        "location": "Bedroom",
        "unit": "lux",
        "status": "OFF"
    }' | jq
echo

# Gate Device
curl -s -X POST "$API_URL/devices" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Gate",
        "type": "gate",
        "location": "Garage",
        "status": "closed"
    }' | jq
echo

echo "Устройства добавлены."

echo "Проверяем список устройств через DMS:"
curl -s "$API_URL/devices" | jq
echo

# ==============================
# Изменение значений
# ==============================

echo "Изменяем значение температуры и статус света..."

# Получаем ID устройств
TEMP_ID=$(curl -s "$API_URL/devices" | jq -r '.[] | select(.name=="Temperature Sensor") | .id')
LIGHT_ID=$(curl -s "$API_URL/devices" | jq -r '.[] | select(.name=="Light") | .id')

echo "TEMP_ID: $TEMP_ID"
echo "LIGHT_ID: $LIGHT_ID"

# Изменяем температуру на 25.0
if [[ -n "$TEMP_ID" ]]; then
  echo "Обновляем температуру до 25.0°C"
  curl -s -X POST "$API_URL/devices/$TEMP_ID/command?isTest=true" \
      -H "Content-Type: application/json" \
      -d '{"value":"25.0"}' | jq
fi

# Включаем свет
if [[ -n "$LIGHT_ID" ]]; then
  echo "Включаем свет"
  curl -s -X POST "$API_URL/devices/$LIGHT_ID/command?isTest=true" \
      -H "Content-Type: application/json" \
      -d '{"status":"ON"}' | jq
fi

echo "Проверяем изменения:"
curl -s "$API_URL/devices" | jq
echo

echo "Проверяем телеметрию для температуры (DeviceID=$TEMP_ID)"
curl -s "$TELEMETRY_URL/telemetry/$TEMP_ID" | jq .

echo "Проверяем телеметрию для света (DeviceID=$LIGHT_ID)"
curl -s "$TELEMETRY_URL/telemetry/$LIGHT_ID" | jq .

# ==============================
# Удаление тестовых устройств
# ==============================

echo "Удаляем тестовые устройства..."

# Gate
GATE_ID=$(curl -s "$API_URL/devices" | jq -r '.[] | select(.name=="Gate") | .id')

for ID in "$TEMP_ID" "$LIGHT_ID" "$GATE_ID"; do
  if [[ -n "$ID" ]]; then
    curl -s -X DELETE "$API_URL/devices/$ID" | jq
  fi
done

echo "Тестовые устройства удалены."

# Проверяем итоговый список
echo "Итоговый список устройств:"
curl -s "$API_URL/devices" | jq
echo
