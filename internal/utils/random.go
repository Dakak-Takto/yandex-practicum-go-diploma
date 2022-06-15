package utils

import "crypto/rand"

func Random(length int) ([]byte, error) {
	var b = make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
