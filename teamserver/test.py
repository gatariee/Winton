import flask 

app = flask.Flask(__name__)

@app.route('/', methods=['POST'])
def index():
    print(flask.request.json)
    return 'OK'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
    