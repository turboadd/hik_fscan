import requests
from config.secret_config import config
import time
from datetime import datetime
import threading
import queue
from config.logger import logger

TELEGRAM_TOKEN = config.get("TelegramBotToken", "")
CHAT_ID = config.get("TelegramChatID", "")

RATE_LIMIT_SECONDS = 1.2

last_sent_times = {}
message_queue = queue.Queue()

def send_telegram_direct(text: str):
    if not TELEGRAM_TOKEN or not CHAT_ID:
        logger.error("[TELEGRAM] Token or Chat ID not configured")
        return
    
    url = f"https://api.telegram.org/bot{TELEGRAM_TOKEN}/sendMessage"

    payload = {
        "chat_id": CHAT_ID,
        "text": text,
        "parse_mode": "HTML"
    }
    try:
        response = requests.post(url, json=payload, timeout=5)
        if response.status_code != 200:
            logger.error(f"[TELEGRAM] Failed: {response.text}")
    except Exception as e:
        logger.error(f"[TELEGRAM] Error: {e}")

def telegram_worker():
    #print("[TELEGRAM] Worker started")
    while True:
        try:
            
            item = message_queue.get()
            
            if item is None:
                logger.error("[TELEGRAM] Worker stopped")
                break
            
            client_id, text = item
            now = datetime.utcnow()
            
            last_time = last_sent_times.get(client_id)
            if last_time:
                elapsed = (now - last_time).total_seconds()
                if elapsed < RATE_LIMIT_SECONDS:
                    logger.error(f"[TELEGRAM] Text from Client_id {client_id} (Not success {RATE_LIMIT_SECONDS}s)")
                    message_queue.task_done()
                    continue

            send_telegram_direct(text)
            last_sent_times[client_id] = now
            time.sleep(1.2)
            message_queue.task_done()
        except Exception as e:
            logger.error(f"[TELEGRAM] Worker error: {e}")
            time.sleep(1)

threading.Thread(target=telegram_worker, daemon=True).start()

def send_telegram_message(client_id: int,text: str):    
    message_queue.put((client_id, text))
    
    
    
