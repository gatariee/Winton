
import base64
import json
import sys
import threading
import tkinter as tk
from tkinter import ttk, scrolledtext, END, font

from Winton.client import Client
from Winton.standalone import get_task_response

from UserInterface.colors import colors
from Utils.print import pretty_print_ls



class AgentTab(ttk.Frame):
    def __init__(self, container, agent_name, **kwargs):
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
        TEAMSERVER = "http://127.0.0.1:80"
        self.client = Client(TEAMSERVER)
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

        else:
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)

    def execute_task(self, command: str):
        task_thread = threading.Thread(target=self.run_task, args=(command,))
        task_thread.start()

    def run_task(self, command):
        match command:
            case "clear":
                self.output_text.delete(1.0, tk.END)

            case "pwd":
                self.output_text.insert(
                    tk.END, f"[*] Tasked beacon to get current working directory\n"
                )
                task_response = get_task_response(self.client, command)
                self.display_task_response(task_response)

            case "ls":
                self.output_text.insert(
                    tk.END, f"[*] Tasked beacon to list directory contents\n"
                )
                task_response = get_task_response(self.client, command)
                files = json.loads(
                    base64.b64decode(task_response["results"][0]["Result"]).decode()
                )
                package = pretty_print_ls(files, self.client)
                self.output_text.insert(tk.END, package)

            case "whoami":
                self.output_text.insert(
                    tk.END, f"[*] Tasked beacon to get current user\n"
                )
                task_response = get_task_response(self.client, command)
                self.display_task_response(task_response)

            case "ps":
                self.output_text.insert(
                    tk.END, f"[*] Tasked beacon to list processes\n"
                )
                task_response = get_task_response(self.client, command)
                self.display_task_response(task_response)

            case "getpid":
                self.output_text.insert(
                    tk.END, f"[*] Tasked beacon to get current process ID\n"
                )
                task_response = get_task_response(self.client, command)
                self.display_task_response(task_response)

            case _:
                self.output_text.insert(tk.END, "[!] Not implemented yet!\n")

    def display_task_response(self, task_response):
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