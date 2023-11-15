package main

import (
	"os"
)

func main() {
	if len(os.Args) != 4 {
		_winton_usage()
		os.Exit(1)
	}

	ts := NewTeamServer(os.Args[1], os.Args[2], os.Args[3])

	_winton_print("Starting Teamserver on [" + ts.IP + ":" + ts.Port + "]")
	// TODO: Separate API for agents and operators, operator will directly interact with the teamserver via the API at 50050
	// Agents will interact with the teamserver via the listener at 80

	_winton_error("Only HTTP listener is currently supported, starting HTTP listener...")

	const HTTP_PORT = "80"
	start_http_listener(ts, HTTP_PORT)
	_winton_print("Default: HTTP Listener started on [" + ts.IP + ":" + HTTP_PORT + "]")

	go ts.checkBeacons()
	select {}
}
