"""

NOT A FUNCTIONAL IMPLANT YET, DO NOT READ THIS CODE PLEASE

"""

import requests
import time

from dataclasses import dataclass

@dataclass
class Beacon:
    UID: str
    Sleep: int = 5

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

def register_agent(data: Agent, URL: str):
    try:
        response = requests.post(URL, json=data)
        if response.status_code == 200:
            print("[*] Agent registered successfully")
            return response.json()
    except Exception as e:
        print(e)

def get_tasks(URL: str):
    try:
        response = requests.get(URL)
        if response.status_code == 200:
            return response.json()
    except Exception as e:
        print(e)

    
if __name__ == "__main__":
    TEAMSERVER = "http://127.0.0.1:50050"
    REGISTER = "/register"
    AGENT_DATA = Agent("127.0.0.1", "localhost", "2", "").to_json()
    AGENT_UID = register_agent(AGENT_DATA, TEAMSERVER + REGISTER)['uid']
    beacon = Beacon(AGENT_UID)
    while True:
        print("[*] ...")
        time.sleep(beacon.Sleep)
        print("[*] Checking for tasks...")
        res = get_tasks(TEAMSERVER + "/tasks/" + AGENT_UID)
        if res:
            print("[*] Tasks found")
            result = f"gatari"
            print("[*] Sending result...")
            
            data = {
                "CommandID": res['tasks'][0]['CommandID'],
                "Result": result
            }
            requests.post(TEAMSERVER + "/results/" + res['tasks'][0]['CommandID'], json=data)