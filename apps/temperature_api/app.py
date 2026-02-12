from flask import Flask, request, jsonify
from datetime import datetime
import random

app = Flask(__name__)

@app.route("/temperature")
def temperature():
    location = request.args.get("location")
    if not location:
        return "location required", 400

    value = round(18 + random.random() * 10, 2)
    sensor_id = str(random.randint(1, 100))

    response = {
        "location": location,
        "value": value,
        "unit": "C",
        "status": "ok",
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "description": "Random temperature",
        "sensor_id": sensor_id,
        "sensor_type": "temperature"
    }
    return jsonify(response)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8088)