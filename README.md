# WintonC2
Yet another Command & Control (C2) server using Golang...

This would be my third time making the exact same barebones C2 but with slight improvements, this time I decided to learn Go while doing it because I was bored.

Very heavily a WIP, C agent coming soon.

![sample](https://i.imgur.com/5owJ9Cg.png)

## Compilation
The `teamserver` compiles to a single golang binary:
```bash
go build .
chmod +x teamserver
./teamserver <ip> <port> <password>
```

The Go agent `Wonton` is compiled to a single golang binary:
```bash
go build .
chmod +x wonton
./wonton
```

The C agent `Wanton` has a MSVC solution file, compile the solution in `Release` mode (not really maintained, probably broken):
```bash
gcc -Wall -std=c99 -o Wanton.exe Commands.c Main.c Utils.c -lcurl
```
_A Makefile is not provided because I am lazy_

## Change Log

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
