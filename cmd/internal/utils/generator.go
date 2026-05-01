package utils

import "crypto/sha256"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func Generate(url string) string {
	hash := sha256.Sum256([]byte(url))
	short := make([]byte, 10)
	for i := range short {
		short[i] = alphabet[int(hash[i])%len(alphabet)]
	}

	return string(short)
}
