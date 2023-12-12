![Winton](./assets/Winton_Logo.png)

<div align="center">
    <h1>Winton</h1>

<i>Yet another Command and Control (C2) framework written in Golang</i>
</div>

Winton is an open-source cross-platform C2 framework written for the purposes of learning adversary emulation and C2 infrastructure.

> 🐒 Winton was designed solely for educational purposes, it is still nowhere close to being operationally functional for red team engagements!

![Cover](./assets/Winton_Banner.png)

## Table of Contents
- [Winton](#winton)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
    - [Teamserver](#teamserver)
    - [Implant](#implant)
    - [Client](#client)
  - [Installation](#installation)
    - [Teamserver](#teamserver-1)
    - [Implant](#implant-1)
    - [Client](#client-1)
  - [Usage](#usage)
    - [Teamserver](#teamserver-2)
    - [Client](#client-2)
  - [OPSEC Considerations](#opsec-considerations--notes)
    - [Implant](#implant-2)
    - [Client](#client-3)
    - [Teamserver](#teamserver-3)

## Features
### Teamserver
> Written in Golang 1.21.1 with Gin (stable on Windows 11 x64/AMD64 & Debian 12.x / Kali 2023.3)
- Support for multiple listeners (HTTP implemented)
- Multiplayer-mode
- Cross-platform binary

### Implant
> Written in Golang 1.21.1 (Windows only*)
- Process migration and process injection
- In-memory .NET assembly execution (creds to: [@ropnop](https://github.com/ropnop/go-clr))
- Built-ins via `os/exec` & `os/user`

### Client
> Dark themed UI written in Python with Tkinter
- Multi-player
- In-memory .NET assembly execution via `execute-assembly`
![execute-assembly](./assets/execute_assembly.png)
  - creds: [SharpAwareness](https://github.com/CodeXTF2/SharpAwareness) by [@CodeXTF2](https://twitter.com/codex_tf2)
  - for some reason, if you try to load .NET assemblies that are too large, the CLR will just not load lol.
- Updated list of supported commands available: [here](./client/Winton/globals.py#)
![Help](./assets/Client_help.png)
- [Athena](https://github.com/gatariee/Athena) - A bot integrated with Winton for collaborative red team operations over Discord 

## Installation
### Winton
```bash
git clone https://github.com/gatariee/Winton
cd Winton
```

### Teamserver
```bash
cd teamserver
make linux # or windows
cd ./bin && chmod +x ./teamserver-x64
```

### Implant
```bash
cd ./implant
make windows
```

### Client
```bash
cd ./client
python3 -m pip install -r requirements.txt
chmod +x ./winton.py
```

## Usage
### Teamserver
```bash
./teamserver-x64 <ip> <port> <password>
```

### Client
```bash
./winton.py
```

## OPSEC Considerations / Notes
### Implant
- The stable implant is written in Go and produces a binary of ~7,747,072 bytes, or ~7.38MB.
- `shell` pipes the input of the operator to `cmd.exe /c {task}`, which spawns a new `cmd.exe` process on the target and returns the output via `stdout` & `stderr`.
- Heavy reliance on Golang's `os/exec` and `os/user` packages for cross-platform compatibility and built-ins (`whoami`, `pwd`, `ls`), may be OPSEC unsafe.
- `inject` uses `CreateRemoteThread` and doesn't check for architecture, may result in the process and/or shellcode crashing- use `ps` to check for architecture before injection.
![Client](./assets/Client_ps.png)
    - `VirtualAllocEx` is called with PAGE_EXECUTE_READWRITE & unbacked memory allocation
    - Thread start address is `0x0`

### Client
- Unencrypted communication with the teamserver over HTTP
- Authentication with teamserver not implemented yet
- Interacts with the listener rather than the teamserver, the operator should be interacting with the internal teamserver API instead of the listener. (modularity)
![Client](./assets/Operator_interaction.png)

### Teamserver
- Unencrypted communication with the implant over HTTP
- Teamserver expects agent to be legitimate and doesn't check for authentication (in fact, the password param used to start the teamserver is completely unused 🤡)
