from Winton.client import Client
from Winton.globals import Teamserver

def dispatch():
    try:
        beacon = Client(Teamserver)
        return beacon.get_agents(Teamserver)
    except Exception as e:
        print("[!] Error: ", e)
        return []