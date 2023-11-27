package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./execute-assembly.exe <path_to_dotnet>")
		os.Exit(1)
	}

	filename := os.Args[1]
	exebytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	runtime.KeepAlive(exebytes)

	runtime := ""
	params := []string{}

	ret2, err := ExecuteByteArray(runtime, exebytes, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] EXE Return Code: %d\n", ret2)
}
