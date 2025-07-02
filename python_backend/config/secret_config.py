import json
import os

CONFIG_PATH = os.path.join(os.path.dirname(__file__), "config.json")

try:
    with open(CONFIG_PATH, "r") as f:
        config = json.load(f)
except Exception as e:
    print("[CONFIG] Failed to load config file:", e)
    config = {}

SECRET = config.get("AuthToken", "")
