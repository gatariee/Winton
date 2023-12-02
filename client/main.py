#!/usr/bin/env python3

from UserInterface.widgets.winton import Winton
from Winton.client import Client
from Utils.config import load

config = load()
Teamserver = config["teamserver"]

def dispatch():
    try:
        beacon = Client(Teamserver)
        return beacon.get_agents(Teamserver)
    except Exception as e:
        print("[!] Error: ", e)
        return []

def main():
    app = Winton(dispatch)
    app.mainloop()


if __name__ == "__main__":
    main()
