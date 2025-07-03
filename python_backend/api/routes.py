from flask import Blueprint, request, jsonify
from auth.auth import verify_jwt_token
# import datetime
from datetime import datetime
from db.init import init_db, SessionLocal
from db.models import Event
from notifier.telegram import send_telegram_message
from config.logger import logger


api_bp = Blueprint("api", __name__)
init_db()

@api_bp.route("/event", methods=["POST"])
def receive_event():
    payload = verify_jwt_token()
    if not payload:
        logger.error("Unauthorized request")
        return jsonify({"error": "Unauthorized"}), 401
    
    data = request.json
    client_id = int(data.get("client_id", 0))
    session = SessionLocal()    

    try:
        event_time = datetime.utcfromtimestamp(int(data.get('time')))
        new_event = Event(
            site_id = data.get("site_id"),
            ip = data.get("ip"),
            device_id = data.get("device_id"),
            client_id = data.get("client_id"),
            event_time = event_time,
        )
        session.add(new_event)
        session.commit()

        msg = f"📡 Event Received\n🔹 Site: {data.get('site_id')}\n📍 Client: {data.get('client_id')}\n🕒 Time: {event_time}"
        send_telegram_message(client_id, msg)

        #print(f"[EVENT RECEIVED] CLIENT ID: {data.get('client_id')}, EVENT TIME: {event_time}")
        return jsonify({"status": "OK"}), 200        
    except Exception as e:
        logger.error(f"[EVENT ERROR ADD] {e}")
        session.rollback()
        return jsonify({"error": str(e)}), 500
    finally:
        session.close()   

@api_bp.route("/health", methods=["GET"])
def health_check():
    return jsonify({
        "status": "healthy",
        "time": datetime.utcnow().isoformat() + "Z"
    }), 200