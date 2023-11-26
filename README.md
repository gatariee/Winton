![Winton](https://i.imgur.com/Pmrv5t7.png)

<div align="center">
<h1>Winton</h1>

<i>Yet another Command and Control (C2) framework written in Go</i>
</div>

I wrote this to learn more about C2 internals, OPSEC considerations in beacon and to learn Golang _(i still can't get function imports to work)_.

![Cover](https://i.imgur.com/xhTM1va.png)

> 🐒 Winton was designed solely for educational purposes, and is still in early stages of development and may be unstable. 

## Features
### Teamserver
> Written in Golang 1.21.1 with Gin (tested on Windows/AMD64)
- Support for multiple listeners- only HTTP is implemented
- Multiple agents & asynchronous callbacks

### Implant
> Written in Golang 1.21.1 (tested on Windows/AMD64) 
- OPSEC Considerations
    - Heavy reliance on Golang's `os/exec` and `os/user` packages for cross-platform compatibility and built-ins (`whoami`, `pwd`, `ls`), may be OPSEC unsafe.
    - `inject` uses `CreateRemoteThread` and doesn't check for architecture, may result in the process and/or shellcode crashing- use `ps` to check for architecture before injection.
        - `ps` uses the `syscall` and `golang.org/x/sys/windows` package to access the WinAPIs, see [source](./implant/Wonton%20(GO)/commands.go#L160)
- Updated list of supported commands available: [here](./client/Winton/globals.py#)
- There are 2 implants available:
    - `Orisa` is written in C and is extremely unstable, and has limited functionality to `ls`, `pwd` and `whoami`.
    - `Sigma` is written in Golang and is much more stable, and has more functionality than `Orisa`.
        - `Sigma` is still in early stages of development, and may be unstable.
        - `Sigma` is the recommended implant to use.

### Client
> Dark themed UI written in Python with Tkinter
- Supports interaction with multiple agents & asynchronous callbacks via multithreading

![Client](https://i.imgur.com/SLLtTob.png)


## Compilation
The `teamserver` compiles to a single golang binary:
```bash
cd ./teamserver && make 
cd ./bin && ./teamserver <ip> <port> <password>
```

The Go agent `Wonton` is compiled to a single golang binary:
```bash
cd './implan/Wonton (GO)\' && make
cd ./bin && ./winton-win64.exe
```

The C agent `Wanton` has a MSVC solution file, compile the solution in `Release` mode (not really maintained, probably broken):
```bash
gcc -Wall -std=c99 -o Wanton.exe Commands.c Main.c Utils.c -lcurl
```
_A Makefile is not provided because I am lazy_