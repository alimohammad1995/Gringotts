package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
)

func TronAddressToHex(tronAddress string) (string, error) {
	decoded, version, err := base58.CheckDecode(tronAddress)

	if err != nil {
		return "", fmt.Errorf("Base58Check decode failed: %v", err)
	}

	if version != 0x41 {
		return "", fmt.Errorf("invalid version byte: expected 0x41, got 0x%X", version)
	}

	if len(decoded) != 20 {
		return "", fmt.Errorf("invalid address length: expected 20 bytes, got %d bytes", len(decoded))
	}

	return "0x" + hex.EncodeToString(decoded), nil
}

func ToByte32(input string) [32]byte {
	var bytes32 [32]byte
	copy(bytes32[:], input)
	return bytes32
}

func FromByte32ToByte(input string) []byte {
	x := ToByte32(input)
	return x[:]
}

func ToByte32SOL(input string) [32]byte {
	var bytes32 [32]byte
	copy(bytes32[:], base58.Decode(input))
	return bytes32
}

func FromByte32ToByteSOL(input string) []byte {
	x := ToByte32SOL(input)
	return x[:]
}
