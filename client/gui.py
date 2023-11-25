import tkinter as tk
import sys
import base64
from tkinter import ttk, scrolledtext, END
from Winton.client import Client
from Utils.print import pretty_print_files, beacon_print, pretty_print, print_agents
from Winton.standalone import gui_get_task_response
from Winton.globals import Tasks

colors = {
    "background": "#1e1e1e",
    "foreground": "#c7c7c7",
    "button_background": "#333333",
    "button_foreground": "#ffffff",
    "text_background": "#252526",
    "text_foreground": "#d4d4d4",
    "listbox_background": "#2d2d2d",
    "listbox_foreground": "#cccccc",
    "select_background": "#3e3e3e",
}

class AgentTab(ttk.Frame):
    def __init__(self, container, agent_name: str, **kwargs):
        super().__init__(container, **kwargs)
        self.agent_name = agent_name
        self.uid = agent_name.split(" | ")[-1]
        self.prompt = f"winton> "
        self.create_widgets()

        self.go()

        
    def go(self):
        TEAMSERVER = "http://127.0.0.1:80"
        self.client = Client(TEAMSERVER=TEAMSERVER)
        self.client.refresh_agents()

        for num, agent in enumerate(self.client.Agent_List, start=1):
            if agent["UID"] == self.uid:
                self.client.choose_agent(num - 1)
                break

        
    def create_widgets(self):
        self.output_text = scrolledtext.ScrolledText(self, bg=colors["text_background"], fg=colors["text_foreground"])
        self.output_text.pack(expand=True, fill=tk.BOTH, padx=10, pady=(5, 0))

        # Command Entry
        self.command_entry = tk.Entry(self, bg=colors["text_background"], fg=colors["text_foreground"], insertbackground=colors["foreground"])
        self.command_entry.insert(0, self.prompt)
        self.command_entry.pack(fill=tk.X, side='bottom', padx=10, pady=(0, 5))
        self.command_entry.bind('<Return>', self.run_command)
        self.command_entry.bind('<Key>', self.on_key_press)

    def handle_task(self, task: str, args: str = ""):
        match task:
            case "pwd":
                self.output_text.insert(tk.END, f"[*] Tasked beacon to get current directory\n")
                self.output_text.update_idletasks()
                task_response = gui_get_task_response(self.client, task.replace(" ", ""))

                response_size = sys.getsizeof(task_response["results"][0]["Result"])
                self.output_text.insert(tk.END, f"[*] {self.client.AgentHostname} called home, sent: {response_size} bytes\n")
                self.output_text.update_idletasks()

                self.output_text.insert(tk.END, f"{base64.b64decode(task_response['results'][0]['Result']).decode()}\n")

            case _:
                self.output_text.insert(tk.END, f"[!] Not implemented yet!\n")

    def run_command(self, event=None):
        command_with_prompt = self.command_entry.get()
        if command_with_prompt.startswith(self.prompt):
            command = command_with_prompt[len(self.prompt):]
            self.output_text.insert(tk.END, f"{command_with_prompt}\n")
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)

            self.handle_task( command )
        else:
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)
    
    def on_key_press(self, event): 
        # budget solution for deleting prompt lol
        if event.keysym in ('BackSpace', 'Delete'):
            if self.command_entry.index(tk.INSERT) <= len(self.prompt):
                return "break"
        elif event.keysym == 'Left':
            if self.command_entry.index(tk.INSERT) <= len(self.prompt):
                return "break"
        elif event.keysym in ('Home', 'Right'):
            self.command_entry.icursor(tk.END)
            return "break"

class WintonApp(tk.Tk):
    def __init__(self, agents):
        super().__init__()
        self.title("WintonC2 Client")
        self.configure(bg=colors["background"])
        self.geometry('1200x600')
        self.agents = agents

        self.style = ttk.Style()
        self.style.theme_use('clam')
        self.style.configure('TNotebook', background=colors["background"], foreground=colors["foreground"])
        self.style.configure('TNotebook.Tab', background=colors["background"], foreground=colors["foreground"], lightcolor=colors["background"], borderwidth=0)
        self.style.map('TNotebook.Tab', background=[('selected', colors["background"])], foreground=[('selected', colors["foreground"])])
        self.style.configure('TFrame', background=colors["background"], foreground=colors["foreground"])

        self.notebook = ttk.Notebook(self, style='TNotebook')
        self.notebook.pack(fill='both', expand=True, padx=10, pady=10)

        self.agent_listbox = tk.Listbox(self, bg=colors["listbox_background"], fg=colors["listbox_foreground"], selectbackground=colors["select_background"], exportselection=False)
        self.agent_listbox.pack(fill=tk.BOTH, side='left', anchor='nw', padx=10, pady=10, expand=True)
        self.agent_listbox.bind('<Double-1>', self.on_agent_double_click)

        self.populate_agents(agents=agents)

    def populate_agents(self, agents):
        for agent in agents:
            self.agent_listbox.insert(END, f"{agent['Hostname']} @ {agent['IP']} | {agent['UID']}")

    def on_agent_double_click(self, event):
        selection = self.agent_listbox.curselection()
        if selection:
            agent_name = self.agent_listbox.get(selection)
            tab_names = [self.notebook.tab(tab, "text") for tab in self.notebook.tabs()]
            if agent_name not in tab_names:
                agent_tab = AgentTab(self.notebook, agent_name)
                self.notebook.add(agent_tab, text=agent_name)
                self.notebook.select(agent_tab)

def main():
    agents = Client.get_agents("http://127.0.0.1:80")
    app = WintonApp(agents)
    app.mainloop()

if __name__ == "__main__":
    main()
