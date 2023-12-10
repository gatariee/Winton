import time
from Winton.client import Client
from Utils.config import load

"""
beacon:
  kill_time: "3600" # stop waiting for result after 1 hour, clear queue and await next task
"""

config = load()
KILL_TIME = config["beacon"]["kill_time"]


def get_task_response(client: Client, task: str, args: str = ""):
    PACKAGE_START = time.time()
    task_request = client.send_task(task + " " + args)

    client.Tasks.append(task_request["uid"])

    time.sleep(int(client.Beacon_Sleep) + 2)

    
    while True:
        if time.time() - PACKAGE_START > int(KILL_TIME):
            print("[!] Beacon died, clearing queue and awaiting next task")
            client.clear_queue()
            return None
        print("[*] Getting result for: ", task_request["uid"])
        ok, res = client.get_results(task=task_request["uid"])
        if ok:
            client.Tasks.remove(task_request["uid"])
            return res
        else:
            time.sleep(int(client.Beacon_Sleep))
            print(res)
            continue
