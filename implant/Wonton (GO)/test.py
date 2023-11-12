from flask import Flask, request
import base64
import json

app = Flask(__name__)

@app.route('/', methods=['GET'])
def index():
    return "OK"

@app.route('/', methods=['POST'])
def post():
    print(request.data)
    data = base64.b64decode(request.data).decode()
    print(data)
    json_data = json.loads(data)
    print(json_data)
    return "OK"


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=50050, debug=True)