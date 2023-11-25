from Winton.types import Command

Tasks: list[Command] = [
    {
        "name": "ls",
        "description": "List files in current directory",
        "usage": "ls"
    },
    {
        "name": "whoami",
        "description": "Get current user",
        "usage": "whoami"
    },
    {
        "name": "pwd",
        "description": "Get current directory",
        "usage": "pwd"
    },
    {
        "name": "ps",
        "description": "List processes",
        "usage": "ps"
    },
    {
        "name": "getpid",
        "description": "Get current process ID",
        "usage": "getpid"
    },
    {
        "name": "shell",
        "description": "Execute shell command",
        "usage": "shell <command>"
    },
    {
        "name": "inject",
        "description": "Inject shellcode",
        "usage": "inject <PID> <path_to_binfile>"
    },
    {
        "name": "help",
        "description": "Display this help menu",
        "usage": "help"
    },
    {
        "name": "clear",
        "description": "Clear the console",
        "usage": "clear"
    },
    {
        "name": "exit",
        "description": "Exit the program",
        "usage": "exit"
    },
    {
        "name": "winton",
        "description": "Winton?",
        "usage": "monkey"
    }
]