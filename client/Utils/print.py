from Winton.client import Client
from Winton.types import File, Agent
from Winton.globals import Tasks
from Utils.winton import WINTON

import sys
import base64
import random

def _beacon_called_home(client: Client, size: bytes):
    print(f"[*] {client.AgentHostname} called home, sent: {size} bytes")

def pretty_print_files(files: list[str, File], client: Client):

    print("\n")

    response_size = sys.getsizeof(files)

    _beacon_called_home(client, response_size)

    max_filename_len = max(len(file["Filename"]) for file in files)

    for file in files:
        file_name = file["Filename"].split("\\")[-1]

        size_in_kb = file["Size"] / 1024
        if file["IsDir"]:
            file_type = "Directory"
        else:
            file_type = "\tFile"

        padding_width = max_filename_len - len(file_name)

        print(f"{file_name}{' ' * padding_width}\t{size_in_kb:.2f}KB\t{file_type}")

    print("\n")


def beacon_print(client: Client, task_response: dict):
    response_size = sys.getsizeof(task_response["results"][0]["Result"])

    _beacon_called_home(client, response_size)

    print(
        f"""
{base64.b64decode(task_response["results"][0]["Result"]).decode()}
    """
    )

def pretty_print(data: str):
    print(f"[*] {data}")


def print_agents(agent_list: list[Agent]):
    if not agent_list:
        print("[!] No agents registered")
        return
    for num, agent in enumerate(agent_list, start=1):
        print(f"{num}. {agent['Hostname']}@{agent['IP']} | {agent['UID']}")
    

def pretty_print_ls(files: list[str, File], client: Client) -> str:
    package = ""
    response_size = sys.getsizeof(files)
    package += f"[*] {client.AgentHostname} called home, sent: {response_size} bytes\n\n"
    max_filename_len = max(len(file["Filename"]) for file in files)
    for file in files:
        file_name = file["Filename"].split("\\")[-1]
        size_in_kb = file["Size"] / 1024
        if file["IsDir"]:
            file_type = "Directory"
        else:
            file_type = "\tFile"
        padding_width = max_filename_len - len(file_name)
        package += f"{file_name}{' ' * padding_width}\t{size_in_kb:.2f}KB\t{file_type}\n"
    
    package += "\n"
    return package

def handle_help__str__(tasks=Tasks):
    package = ""
    header = "Command             Description"
    separator = "-------             -----------"

    package += "\nWintonC2 Commands\n================\n"
    package += "\n"
    package += header
    package += "\n"
    package += separator
    package += "\n"

    for command in tasks:
        name_space = ' ' * (20 - len(command['name']))
        
        line = f"{command['name']}{name_space}{command['description']}"
        package += line
        package += "\n"

    package += "\n"

    return package

def handle_winton():
    package = ""
    package += "\n"
    package += "Winton\n"
    package += "=====\n"
    package += "\n"

    package += "Winton is a Command and Control (C2) framework written in Golang by @gatari"
    package += "\n"

    package += "https://github.com/gatariee/Winton"

    package += "\n"

    package += random.choice(WINTON)

    package += "\n"

    return package