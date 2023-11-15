package main

import (
	"fmt"
)

func _winton_print(data string) {
	fmt.Println("[*] " + data)
}

func _winton_error(data string) {
	fmt.Println("[!] " + data)
}

func _winton_usage() {
	fmt.Println("Usage: ./teamserver <ip> <port> <password>")
}
