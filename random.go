package core

import "crypto/rand"

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
const digits = "0123456789"

func GenerateRandom(numberOnly bool, length uint8) (string, error) {
	var template string
	if numberOnly {
		template = digits
	} else {
		template = letters
	}
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = template[b%byte(len(template))]
	}
	return string(bytes), nil
}

func GenerateRandomBytes(length uint8) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
