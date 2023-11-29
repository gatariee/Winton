from Winton.types import Command

Tasks: list[Command] = [
    {"name": "ls", "description": "List files in current directory", "usage": "ls"},
    {"name": "whoami", "description": "Get current user", "usage": "whoami"},
    {"name": "cat", "description": "Read file", "usage": "cat <path_to_file>"},
    {"name": "pwd", "description": "Get current directory", "usage": "pwd"},
    {"name": "ps", "description": "List processes", "usage": "ps"},
    {"name": "getpid", "description": "Get current process ID", "usage": "getpid"},
    {"name": "help", "description": "Display this help menu", "usage": "help"},
    {"name": "exit", "description": "Exit the program", "usage": "exit"},
    {"name": "winton", "description": "Winton?", "usage": "monkey"},
    {
        "name": "shell",
        "description": "Execute shell command",
        "usage": "shell <command>",
    },
    {
        "name": "inject",
        "description": "Inject shellcode into a process",
        "usage": "inject <PID> <path_to_binfile>",
    },
    {
        "name": "execute-assembly",
        "description": "Execute .NET assembly in memory",
        "usage": "execute-assembly <path_to_assembly>",
    }
]

Teamserver = "http://127.0.0.1:80"
