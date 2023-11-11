package main

import (
	"fmt"
	"encoding/json"
)

func main() {
	files, err := ls("C:\\Users\\PC\\Desktop\\git\\WintonC2")
	if err != nil {
		fmt.Println(err)
		return

	}

	jsonFiles, err := json.Marshal(files)
	if err != nil {
		fmt.Println(err)
		return
	}

	enc := b64_encode(jsonFiles)
	fmt.Println(enc)
	
	dec, err:= b64_decode(enc)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(dec))
}
