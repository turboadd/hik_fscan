from flask import Blueprint, request, jsonify
from auth.auth import verify_jwt_token
# import datetime
from datetime import datetime

api_bp = Blueprint("api", __name__)

@api_bp.route("/event", methods=["POST"])
def receive_event():
    payload = verify_jwt_token()
    if not payload:
        return jsonify({"error": "Unauthorized"}), 401
    
    data = request.json

    print(f"[EVENT RECEIVED] USER ID: {data.get('user_id')}, EVENT TIME: {datetime.utcfromtimestamp(int(data.get('time')))}")
    return jsonify({"status": "OK"}), 200

@api_bp.route("/health", methods=["GET"])
def health_check():
    return jsonify({
        "status": "healthy",
        "time": datetime.utcnow().isoformat() + "Z"
    }), 200