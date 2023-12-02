package utils

import (
	"encoding/base64"
)

func Base64_Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64_Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
