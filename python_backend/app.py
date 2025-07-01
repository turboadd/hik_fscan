from flask import Flask, request
app = Flask(__name__)

@app.route("/event", methods=["POST"])
def receive_event():
    print("[EVENT RECEIVED]", request.json)
    return "OK", 200

if __name__ == "__main__":
    app.run(port=5000, debug=True)