import requests
import time
from dataclasses import dataclass

URL = "http://127.0.0.1/"

@dataclass
class Agent:
    IP: str
    Hostname: str
    Sleep: str
    UID: str

    def to_json(self):
        return {
            "IP": self.IP,
            "Hostname": self.Hostname,
            "Sleep": self.Sleep,
            "UID": self.UID
        }

@dataclass
class Task:
    UID: str
    Command: str

    def to_json(self):
        return {
            "UID": self.UID,
            "Command": self.Command
        }
    
@dataclass
class Result:
    CommandID: str
    Result: str

    def to_json(self):
        return {
            "CommandID": self.CommandID,
            "Result": self.Result
        }



def send_get(url: str):
    try:
        response = requests.get(url)
        print(response.text)
    except Exception as e:
        print(e)

def register_agent(data: Agent):
    try:
        API = "register"
        response = requests.post(URL + API, json=data)
        print(response.text)
    except Exception as e:
        print(e)

def send_task(data: Task):
    try:
        API = "tasks"
        command = { "Command": data['Command'], "CommandID": ""}
        response = requests.post(URL + API + "/" + data['UID'], json=command)
        print(response.text)
    except Exception as e:
        print(e)

def get_tasks(id: str):
    try:
        API = "tasks"
        response = requests.get(URL + API + "/" + id)
        print(response.text)
    except Exception as e:
        print(e)

def send_result(data: Result):
    try:
        API = "results"
        response = requests.post(URL + API + "/" + data['CommandID'], json=data)
        print(response.text)
    except Exception as e:
        print(e)

def get_results(id: str):
    try:
        API = "results"
        response = requests.get(URL + API + "/" + id)
        print(response.text)
    except Exception as e:
        print(e)

def run_tests():
    data = Agent("127.0.0.1", "test", "5", "test").to_json()
    register_agent(data)
    time.sleep(1)

if __name__ == "__main__":
    # data = Agent("127.0.0.1", "test", "5", "test").to_json()
    # register_agent(data)

    # data = Task("3ud9FuATeW", "ayaya2").to_json()
    # send_task(data)

    # get_tasks("3ud9FuATeW")

    data = Result("MLcAfGWUIK", "Test").to_json()
    send_result(data)

