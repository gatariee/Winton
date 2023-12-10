import requests
from dataclasses import dataclass
from Winton.types import Agent, CommandData

@dataclass
class Client:
    Agent_List: list[Agent]
    Tasks: list[str]
    AgentID: str = ""
    AgentHostname: str = ""
    AgentIP: str = ""
    Beacon_Sleep: str = ""
    Teamserver: str = ""

    def __init__(self, TEAMSERVER: str):
        self.Agent_List = self.get_agents(TEAMSERVER)
        self.Teamserver = TEAMSERVER
        self.Tasks = []

    @classmethod
    def get_agents(cls, URL: str) -> list[Agent]:
        try:
            URL = URL + "/agents"
            response = requests.get(URL)
            if response.status_code == 200:
                return response.json()["agents"]
        except Exception as e:
            print(e)
            return []

    def refresh_agents(self):
        self.Agent_List = self.get_agents(self.Teamserver)

    def send_task(self, task: str):
        try:
            URL = self.Teamserver + "/tasks/" + self.AgentID
            command_data = CommandData(CommandID="", Command=task)
            data = command_data.winton()
            response = requests.post(URL, json=data)
            if response.status_code == 200:
                return response.json()
        except Exception as e:
            print(e)

    def get_results(self, task) -> (bool, list[dict]):
        URL = self.Teamserver + "/results/" + task
        try:
            response = requests.get(URL) # this can error
            if response.status_code == 200:
                return True, response.json()["results"]
            
            if response.json()["message"] == "No results found":
                return False, response.json()["message"]
        
        except Exception as e:
            # if an error is thrown, let the caller handle it
            return False, e

    def display_agents(self):
        if len(self.Agent_List) == 0:
            print("[!] No agents registered")
            return

        for num, agent in enumerate(self.Agent_List, start=1):
            print(f"{num}. {agent['Hostname']}@{agent['IP']} | {agent['UID']}")

    def choose_agent(self, num: int):
        self.AgentID = self.Agent_List[num]["UID"]
        self.Beacon_Sleep = self.Agent_List[num]["Sleep"]
        self.AgentHostname = self.Agent_List[num]["Hostname"]
        self.AgentIP = self.Agent_List[num]["IP"]

    def reset_agent(self):
        self.AgentID = ""
        self.Beacon_Sleep = ""
        self.AgentHostname = ""
        self.AgentIP = ""