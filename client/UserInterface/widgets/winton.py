import tkinter as tk
from tkinter import ttk, END, font

from UserInterface.globals import colors
from UserInterface.widgets.agent import AgentTab

from Winton.types import Agent

class Winton(tk.Tk):
    def __init__(self, fetch_agents, **kwargs):
        super().__init__()
        self.title("Winton")
        self.configure(bg=colors["background"])
        self.geometry("1400x800")
        self.fetch_agents = fetch_agents
        self.agents = self.fetch_agents()
        
        self.setup_style()
        self.setup_notebook()
        self.populate_agents(self.agents)
        self.schedule_agent_update()
    


    def setup_style(self):
        modern_font = font.nametofont("TkDefaultFont")
        modern_font.configure(family="Consolas", size=12)

        style = ttk.Style()
        style.theme_use("clam")

        style.configure(
            "TNotebook",
            background=colors["background"],
            foreground=colors["foreground"],
        )
        style.configure(
            "TNotebook.Tab",
            background=colors["background"],
            foreground=colors["foreground"],
            lightcolor=colors["background"],
            borderwidth=0,
        )
        style.map(
            "TNotebook.Tab",
            background=[("selected", colors["background"])],
            foreground=[("selected", colors["foreground"])],
        )
        style.configure(
            "TFrame", background=colors["background"], foreground=colors["foreground"]
        )

        style.configure(
            "TButton",
            background=colors["button_background"],
            foreground=colors["button_foreground"],
            font=modern_font,
            borderwidth=0,
        )
        style.map("TButton", background=[("active", colors["button_background"])])

        style.configure(
            "TEntry",
            background=colors["text_background"],
            foreground=colors["text_foreground"],
            fieldbackground=colors["text_background"],
            borderwidth=0,
        )
        style.configure(
            "TListbox",
            background=colors["listbox_background"],
            foreground=colors["listbox_foreground"],
            selectbackground=colors["select_background"],
            borderwidth=0,
        )
        style.configure(
            "TLabel",
            background=colors["background"],
            foreground=colors["foreground"],
            font=modern_font,
        )

    def setup_notebook(self):
        self.notebook = ttk.Notebook(self, style="TNotebook")
        self.notebook.pack(fill="both", expand=True, padx=10, pady=10)

        self.agent_listbox = tk.Listbox(
            self,
            bg=colors["listbox_background"],
            fg=colors["listbox_foreground"],
            selectbackground=colors["select_background"],
            exportselection=False,
        )
        self.agent_listbox.pack(
            fill=tk.BOTH, side="left", anchor="nw", padx=10, pady=10, expand=True
        )
        self.agent_listbox.bind("<Double-1>", self.on_agent_double_click)

    def populate_agents(self, agents: list[Agent]):
        self.agent_listbox.delete(0, END) 
        for agent in agents:
            self.agent_listbox.insert(
                END, f"{agent['Hostname']} @ {agent['IP']} | {agent['UID']}"
            )
        
    def schedule_agent_update(self):
        self.agents = self.fetch_agents()
        self.populate_agents(self.agents)
        self.after(5000, self.schedule_agent_update) 


    def on_agent_double_click(self, event):
        selection = self.agent_listbox.curselection()
        if selection:
            agent_name = self.agent_listbox.get(selection)
            self.open_agent_tab(agent_name)

    def open_agent_tab(self, agent_name: str):
        tab_names = [self.notebook.tab(tab, "text") for tab in self.notebook.tabs()]
        if agent_name not in tab_names:
            agent_tab = AgentTab(self.notebook, agent_name)
            self.notebook.add(agent_tab, text=agent_name)
            self.notebook.select(agent_tab)