package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateID() (int, error) {
	maxInt32 := big.NewInt(1<<31 - 1)
	n, err := rand.Int(rand.Reader, maxInt32)
	if err != nil {
		return -1, err
	}

	return int(n.Int64()), nil
}
