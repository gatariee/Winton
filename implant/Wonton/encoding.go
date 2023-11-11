package main

import (
	"encoding/base64"
)

func b64_encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func b64_decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
