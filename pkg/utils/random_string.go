package utils

import (
	"crypto/rand"
	"math/big"
)

func RandomString() (string, error) {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	s := ""
	for i := 0; i < 32; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))-1))
		if err != nil {
			return "", err
		}
		s += string(letters[n.Int64()])
	}
	return s, nil
}
