from flask import Flask
from config.secret_config import config
from api.routes import api_bp


app = Flask(__name__)

app.register_blueprint(api_bp)

if __name__ == "__main__":
    #print(jwt.__file__)
    app.run(port=config.get("listen_port",""), debug=True)