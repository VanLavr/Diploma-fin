package generator

import (
	"crypto/rand"
	"math/big"
)

const DEFAULTLEN = 12

func GeneratePassword(length int) (string, error) {
	const (
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits       = "0123456789"
		symbols      = "!@#$%^&*()-_=+,.?/:;{}[]~"
	)

	allChars := lowerLetters + upperLetters + digits + symbols

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[num.Int64()]
	}

	return string(password), nil
}
