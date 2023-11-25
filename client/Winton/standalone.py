import time
from Winton.client import Client

def get_task_response(client: Client, task: str, args: str = ""):
    task_request = client.send_task(task + " " + args)
    client.Tasks = task_request["uid"]
    time.sleep(int(client.Beacon_Sleep) + 2)
    return client.get_results()

def gui_get_task_response(client: Client, task: str, args: str = ""):
    task_request = client.send_task(task + " " + args)
    client.Tasks = task_request["uid"]
    time.sleep(int(client.Beacon_Sleep) + 2)
    return client.get_results()