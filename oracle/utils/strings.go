package utils

import "encoding/hex"

func ToHex(input []byte) string {
	return "0x" + hex.EncodeToString(input)
}

func FromHex(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}
