from Winton.client import Client
from UserInterface.widgets.winton import Winton

def main():
    try:
        agents = Client.get_agents("http://127.0.0.1:80")
    except Exception as e:
        print(e)
        agents = []
        
    app = Winton(agents)
    app.mainloop()


if __name__ == "__main__":
    main()
