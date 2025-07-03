import jwt
from jwt.exceptions import ExpiredSignatureError, InvalidTokenError
from flask import request
from config.secret_config import config
from config.logger import logger

def verify_jwt_token():
    auth_header = request.headers.get("Authorization")
    if not auth_header or not auth_header.startswith("Bearer "):
        return None
    token = auth_header.split(" ")[1]
    #print(token)

    try:
        payload = jwt.decode(token, config["AuthToken"], algorithms=["HS256"])
        return payload
    except jwt.ExpiredSignatureError:
        logger.warning("[AUTH] Token expired")
        return None
    except jwt.InvalidTokenError:
        logger.error("[AUTH] Invalid token")
        return None