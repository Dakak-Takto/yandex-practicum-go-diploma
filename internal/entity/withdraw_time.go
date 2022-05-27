package entity

import (
	"fmt"
	"time"
)

type withdrawTime time.Time

func (uploadedAt *withdrawTime) MarshalJSON() ([]byte, error) {

	t := time.Time(*uploadedAt).Format(time.RFC3339)
	result := fmt.Sprintf("\"%s\"", t)

	return []byte(result), nil
}
