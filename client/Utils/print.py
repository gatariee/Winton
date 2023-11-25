from Winton.client import Client
from Winton.types import File, Agent
import sys
import base64

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