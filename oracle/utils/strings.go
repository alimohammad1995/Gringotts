package utils

import (
	"encoding/hex"
	"strings"
)

func ToHex(input []byte) string {
	return "0x" + hex.EncodeToString(input)
}

func FromHex(input string) []byte {
	if strings.HasPrefix(input, "0x") {
		input = input[2:]
	}
	res, _ := hex.DecodeString(input)
	return res
}

func ZeroOneStringToByteArray(bitString string) []byte {
	size := (len(bitString) + 7) / 8
	byteArray := make([]byte, size)

	for i := 0; i < size; i++ {
		var b byte
		for j := 0; j < 8; j++ {
			bitIndex := i*8 + j
			if bitIndex < len(bitString) && bitString[bitIndex] == '1' {
				b |= 1 << (7 - j)
			}
		}

		byteArray[i] = b
	}
	return byteArray
}

func SplitStringIntoChunks(s string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}

	return chunks
}
