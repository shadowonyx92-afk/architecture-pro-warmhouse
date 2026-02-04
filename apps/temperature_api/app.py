from flask import Flask, request, jsonify
from datetime import datetime
import random

app = Flask(__name__)
sensorType = {
    "Living Room": "1",
    "Bedroom": "2",
    "Kitchen": "3"
}

@app.route("/temperature")
def temperature():
    location = request.args.get("location")
    if not location:
        return "location required", 400

    value = round(18 + random.random() * 10, 2)
    response = {
        "location": location,
        "value": value,
        "unit": "C",
        "status": "ok",
        "timestamp": datetime.now(),
        "description": "Random temperature",
        "sensor_id": sensorType.get(value, lambda: "0"),
        "sensor_type": "temperature"
    }
    return jsonify(response)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8081)