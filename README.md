# WintonC2
Yet another Command & Control (C2) server using Golang...

This would be my third time making the exact same barebones C2 but with slight improvements, this time I decided to learn Go while doing it because I was bored.

Very heavily a WIP, C agent coming soon.

## Compilation
Only the `teamserver` needs to be compiled,
```bash
go build .
chmod +x teamserver
./teamserver <ip> <port> <password>
```
