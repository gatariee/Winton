/*
HTTP Implant
*/

package main

import (
	"fmt"
	"Winton/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Println("[!] Something went very wrong: ", err)

	}
}

