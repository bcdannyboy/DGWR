package utils

import (
	"encoding/binary"
	"math/rand"
)

// CryptoRandFloat64 generates a cryptographically secure random float64 between 0.0 and 1.0
func CryptoRandFloat64() (float64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	// Convert the bytes to a uint64, then to a float64 between 0 and 1
	return float64(binary.LittleEndian.Uint64(b[:])) / (1 << 64), nil
}

func CoinFlip() (bool, error) {
	rf, err := CryptoRandFloat64()
	if err != nil {
		return false, err
	}

	return rf > 0.5, nil
}
