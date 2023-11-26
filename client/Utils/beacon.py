from Winton.client import Client
from Winton.globals import Teamserver

def dispatch():
    beacon = Client(Teamserver)
    return beacon.get_agents(Teamserver)