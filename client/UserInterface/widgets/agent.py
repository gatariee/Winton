import base64
import json
import sys
import threading
import tkinter as tk
from tkinter import ttk, scrolledtext, END, font

from Winton.client import Client
from Winton.standalone import get_task_response
from Winton.globals import Teamserver
from Winton.types import ResultList

from UserInterface.globals import colors
from Utils.print import pretty_print_ls, handle_help__str__, handle_winton


class AgentTab(ttk.Frame):
    def __init__(self, container: ttk.Notebook, agent_name: str, **kwargs):
        super().__init__(container, **kwargs)
        self.agent_name = agent_name
        self.uid = agent_name.split(" | ")[-1]
        self.prompt = f"winton>> "
        self.setup_style()
        self.initialize_client()
        self.create_widgets()

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

    def initialize_client(self):
        self.client = Client(Teamserver)
        self.client.refresh_agents()
        for num, agent in enumerate(self.client.Agent_List, start=1):
            if agent["UID"] == self.uid:
                self.client.choose_agent(num - 1)
                break

    def create_widgets(self):
        self.output_text = scrolledtext.Text(
            self, bg=colors["text_background"], fg=colors["text_foreground"]
        )
        self.output_text.pack(expand=True, fill=tk.BOTH, padx=0, pady=(0, 10))
        self.output_text.configure(font=("Consolas", 12))
        self.output_text.insert(
            tk.END,
            f"[*] Beacon {self.client.AgentHostname} @ {self.client.AgentIP} connected\n",
        )

        self.command_entry = tk.Entry(
            self,
            bg=colors["text_background"],
            fg=colors["text_foreground"],
            insertbackground=colors["foreground"],
            font=("Consolas", 12),
        )
        self.command_entry.insert(0, self.prompt)
        self.command_entry.pack(fill=tk.X, side="bottom", padx=10, pady=(0, 5))
        self.command_entry.bind("<Return>", self.run_command)
        self.command_entry.bind("<Key>", self.prevent_prompt_deletion)

        self.command_entry.bind("<Up>", self.prev_command)
        self.command_entry.bind("<Down>", self.next_command)
        self.command_history = []
        self.history_index = 0
    
    def prev_command(self, event):
        if self.command_history and self.history_index > 0:
            self.history_index -= 1
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt + self.command_history[self.history_index])
            self.command_entry.icursor(tk.END)
            return "break"

    def next_command(self, event):
        if self.command_history and self.history_index < len(self.command_history) - 1:
            self.history_index += 1
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt + self.command_history[self.history_index])
            self.command_entry.icursor(tk.END)
            return "break"
        elif self.command_history and self.history_index == len(self.command_history) - 1:
            self.history_index += 1
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)
            self.command_entry.icursor(tk.END)
            return "break"
    

    def prevent_prompt_deletion(self, event):
        if event.keysym in ("BackSpace", "Delete", "Left") and self.command_entry.index(
            tk.INSERT
        ) <= len(self.prompt):
            return "break"
        elif event.keysym in ("Home", "Right"):
            self.command_entry.icursor(tk.END)
            return "break"

    def run_command(self, event=None):
        command_with_prompt = self.command_entry.get()
        if command_with_prompt.startswith(self.prompt):
            command = command_with_prompt[len(self.prompt) :]
            self.output_text.insert(tk.END, f"{command_with_prompt}\n")
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)
            self.execute_task(command)
            self.command_history.append(command)
            self.history_index = len(self.command_history)

        else:
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)

    def scroll_to_end(self):
        self.output_text.see(tk.END)

    def execute_task(self, command: str):
        task_thread = threading.Thread(target=self.run_task, args=(command,))
        task_thread.start()

    def handle_shell(self, command: str):
        self.output_text.insert(
            tk.END, f"[*] Tasked beacon to execute shell command: {command}\n"
        )
        self.scroll_to_end()
        task_response = get_task_response(
            self.client, "shell", " ".join(command.split(" ")[1:])
        )
        self.display_task_response(task_response)

    def run_task(self, command: str):
        if command.startswith("shell"):
            self.handle_shell(command)
            return
        
        match command:
            case "clear":
                self.handle_clear()
            case "help":
                self.handle_help()
            case "winton":
                self.handle_winton()
            case "pwd":
                self.handle_pwd()
            case "ls":
                self.handle_ls()
            case "whoami":
                self.handle_whoami()
            case "ps":
                self.handle_ps()
            case "getpid":
                self.handle_getpid()
            case _:
                self.handle_default()

        self.scroll_to_end()

    def handle_clear(self):
        self.output_text.delete(1.0, tk.END)

    def handle_help(self):
        package = handle_help__str__()
        self.output_text.insert(tk.END, package)

    def handle_winton(self):
        package = handle_winton()
        self.output_text.insert(tk.END, package)

    def handle_pwd(self):
        self.generic_task_handler("pwd", "get current working directory")

    def handle_ls(self):
        self.output_text.insert(tk.END, f"[*] Tasked beacon to list files in .\n")
        task_response = get_task_response(self.client, "ls")
        files = json.loads(
            base64.b64decode(task_response["results"][0]["Result"]).decode()
        )
        package = pretty_print_ls(files, self.client)
        self.output_text.insert(tk.END, package)

    def handle_whoami(self):
        self.generic_task_handler("whoami", "get current user")

    def handle_ps(self):
        self.generic_task_handler("ps", "get processes")

    def handle_getpid(self):
        self.generic_task_handler("getpid", "getcurrent process ID")

    def handle_default(self):
        self.output_text.insert(tk.END, "[!] Not implemented yet!\n")

    def generic_task_handler(self, command: str, task_description: str):
        self.output_text.insert(tk.END, f"[*] Tasked beacon to {task_description}\n")
        task_response = get_task_response(self.client, command)
        self.display_task_response(task_response)


    def display_task_response(self, task_response: ResultList):
        print(task_response)
        if "results" in task_response and len(task_response["results"]) > 0:
            response_size = sys.getsizeof(task_response["results"][0]["Result"])
            self.output_text.insert(
                tk.END,
                f"[*] {self.client.AgentHostname} called home, sent: {response_size} bytes\n",
            )
            result = base64.b64decode(task_response["results"][0]["Result"]).decode()
            self.output_text.insert(tk.END, f"\n{result}\n\n")
        else:
            self.output_text.insert(tk.END, "[!] No response received\n")