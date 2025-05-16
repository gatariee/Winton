import base64
import json
import sys
import threading
import tkinter as tk
import glob

from tkinter import ttk, scrolledtext, END, font

from Winton.client import Client
from Winton.standalone import get_task_response
from Winton.types import ResultList
from Winton.globals import Tasks

from UserInterface.globals import colors
from Utils.print import pretty_print_ls, handle_help__str__, handle_winton, handle_usage
from Utils.config import load

config = load()
ip = config["teamserver"]["ip"]
port = config["teamserver"]["port"]
Teamserver = f"http://{ip}:{port}"


class AgentTab(ttk.Frame):

    def __init__(self, container: ttk.Notebook, agent_name: str, **kwargs):
        super().__init__(container, **kwargs)
        self.uid = agent_name.split("[")[1].split("]")[0]  # lol this is so bad
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

        style.configure("TFrame",
                        background=colors["background"],
                        foreground=colors["foreground"])

        style.configure(
            "TButton",
            background=colors["button_background"],
            foreground=colors["button_foreground"],
            font=modern_font,
            borderwidth=0,
        )
        style.map("TButton",
                  background=[("active", colors["button_background"])])

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
        self.output_text = scrolledtext.Text(self,
                                             bg=colors["text_background"],
                                             fg=colors["text_foreground"])
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

        self.tasks = [task["name"] for task in Tasks]
        self.command_entry.bind("<Tab>", self.tab_complete)

    def tab_complete(self, event):
        current_text = self.command_entry.get()

        if "cat" in current_text or "execute-assembly" in current_text:
            parts = current_text.split(" ")
            if len(parts) > 1:
                file_path_fragment = parts[-1]
                matches = glob.glob(file_path_fragment + "*")
                if len(matches) == 1:
                    self.command_entry.delete(0, tk.END)
                    self.command_entry.insert(
                        0, " ".join(parts[:-1] + [matches[0]]))
                    self.command_entry.icursor(tk.END)
                    return "break"
                elif len(matches) > 1:
                    self.output_text.insert(tk.END, "\n")
                    for match in matches:
                        self.output_text.insert(tk.END, f"{match}\n")
                    self.output_text.insert(tk.END, "\n")
                    self.scroll_to_end()
                    return "break"

        if current_text.startswith(self.prompt):
            current_text = current_text[len(self.prompt):]
        if current_text == "":
            return "break"

        matching_tasks = [
            task for task in self.tasks if task.startswith(current_text)
        ]
        if len(matching_tasks) == 1:
            self.command_entry.delete(len(self.prompt), tk.END)
            self.command_entry.insert(len(self.prompt), matching_tasks[0])
            self.command_entry.icursor(tk.END)
            return "break"
        elif len(matching_tasks) > 1:
            self.output_text.insert(tk.END, "\n")
            for task in matching_tasks:
                self.output_text.insert(tk.END, f"{task}\n")
            self.output_text.insert(tk.END, "\n")
            self.scroll_to_end()
            return "break"
        else:
            return "break"

    def prev_command(self, event):
        if self.command_history and self.history_index > 0:
            self.history_index -= 1
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(
                0, self.prompt + self.command_history[self.history_index])
            self.command_entry.icursor(tk.END)
            return "break"

    def next_command(self, event):
        if self.command_history and self.history_index < len(
                self.command_history) - 1:
            self.history_index += 1
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(
                0, self.prompt + self.command_history[self.history_index])
            self.command_entry.icursor(tk.END)
            return "break"
        elif (self.command_history
              and self.history_index == len(self.command_history) - 1):
            self.history_index += 1
            self.command_entry.delete(0, tk.END)
            self.command_entry.insert(0, self.prompt)
            self.command_entry.icursor(tk.END)
            return "break"

    def prevent_prompt_deletion(self, event):
        if event.keysym in ("BackSpace", "Delete",
                            "Left") and self.command_entry.index(
                                tk.INSERT) <= len(self.prompt):
            return "break"
        elif event.keysym in ("Home", "Right"):
            self.command_entry.icursor(tk.END)
            return "break"

    def run_command(self, event=None):
        command_with_prompt = self.command_entry.get()
        if command_with_prompt.startswith(self.prompt):
            command = command_with_prompt[len(self.prompt):]
            self.output_text.insert(tk.END, f"{command_with_prompt}\n")
            self.scroll_to_end()
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
        task_thread = threading.Thread(target=self.run_task, args=(command, ))
        task_thread.start()

    def run_task(self, command: str):
        if command.startswith("shell"):
            self.handle_shell(command)
            return

        if command.startswith("cat"):
            self.handle_cat(command)
            return

        if command.startswith("execute-assembly"):
            self.handle_execute_assembly(command)
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

    def handle_execute_assembly(self, command: str):
        if len(command.split(" ")) < 2:
            usage = handle_usage("execute-assembly")
            self.output_text.insert(tk.END, usage)
            return
        file_name = command.split(" ")[1]
        print(f"debug: {file_name}")
        args = command.split(" ")[2:]
        self.output_text.insert(
            tk.END, f"[*] Tasked beacon to run .NET program: {file_name} with args: {args}\n")

        try:
            with open(file_name, "rb") as f:
                assembly_bytes = f.read()
                encoded_bytes = base64.b64encode(assembly_bytes).decode()
        except FileNotFoundError:
            self.output_text.insert(tk.END,
                                    f"[!] File not found: {file_name}\n")
            return

        task_response = get_task_response(self.client, "execute-assembly",
                                          f"{encoded_bytes} "+ " ".join(args))
        self.display_task_response(task_response)

    def handle_shell(self, command: str):
        self.output_text.insert(
            tk.END, f"[*] Tasked beacon to execute shell command: {command}\n")
        self.scroll_to_end()
        task_response = get_task_response(self.client, "shell",
                                          " ".join(command.split(" ")[1:]))
        self.display_task_response(task_response)

    def handle_cat(self, command: str):
        self.output_text.insert(
            tk.END,
            f"[*] Tasked beacon to read contents of file: {command[4:]}\n")
        self.scroll_to_end()
        task_response = get_task_response(self.client, "cat",
                                          " ".join(command.split(" ")[1:]))
        self.display_task_response(task_response)

    def handle_ls(self):
        self.output_text.insert(tk.END,
                                f"[*] Tasked beacon to list files in .\n")
        task_response = get_task_response(self.client, "ls")
        files = json.loads(
            base64.b64decode(task_response[0]["Result"]).decode())
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
        self.output_text.insert(tk.END,
                                f"[*] Tasked beacon to {task_description}\n")
        task_response = get_task_response(self.client, command)
        print(task_response)
        self.display_task_response(task_response)

    def display_task_response(self, task_response: ResultList):
        response_size = sys.getsizeof(task_response[0]["Result"])
        self.output_text.insert(
            tk.END,
            f"[*] {self.client.AgentHostname} called home, sent: {response_size} bytes\n",
        )
        result = base64.b64decode(task_response[0]["Result"]).decode()
        self.output_text.insert(tk.END, f"\n{result}\n\n")
