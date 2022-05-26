package entity

import (
	"errors"
)

type orderStatus uint8

const (
	orderStatusNew orderStatus = 0 + iota
	orderStatusProcessing
	orderStatusInvalid
	orderStatusProcessed
)

func (status *orderStatus) MarshalJSON() ([]byte, error) {

	switch *status {

	case orderStatusNew:
		return []byte("\"new\""), nil

	case orderStatusProcessing:
		return []byte("\"processing\""), nil

	case orderStatusInvalid:
		return []byte("\"invalid\""), nil

	case orderStatusProcessed:
		return []byte("\"processed\""), nil

	}

	return nil, errors.New("unknown order status")
}
