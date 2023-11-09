# WintonC2
Yet another Command & Control (C2) server using Golang...

This would be my third time making the exact same barebones C2 but with slight improvements, this time I decided to learn Go while doing it because I was bored.

Very heavily a WIP, C agent coming soon.

![sample](https://i.imgur.com/5owJ9Cg.png)

## Compilation
Only the `teamserver` needs to be compiled,
```bash
go build .
chmod +x teamserver
./teamserver <ip> <port> <password>
```

## Change Log

#### 9/11/2023 - Beacons can now go offline!
- If the callback from the last beacon is over `Agent.Sleep + 5` seconds, the beacon will be marked as down (_but should still listening for callbacks_).
   - In its current state, the agent gets completely removed from the `AgentList[]` and `AgentCallbacks[]` if it goes offline, this will be changed in the future to allow for offline agents to still be in the list.
    ![9/11/2023](https://i.imgur.com/CZm1eGe.png)
- The change of beacon state is also reflected in the `AgentList[]` and `AgentCallbacks[]` tables
    ![9/11/2023](https://i.imgur.com/p87EHej.png)