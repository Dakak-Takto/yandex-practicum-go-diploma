package utils

import "crypto/rand"

func Random(lenght int) []byte {
	var b = make([]byte, lenght)

	_, err := rand.Read(b)
	if err != nil {
		return nil
	}

	return b
}
