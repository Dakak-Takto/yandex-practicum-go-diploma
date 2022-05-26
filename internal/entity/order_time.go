package entity

import (
	"fmt"
	"time"
)

type orderTime time.Time

func (uploadedAt *orderTime) MarshalJSON() ([]byte, error) {

	t := time.Time(*uploadedAt).Format(time.RFC3339)
	result := fmt.Sprintf("\"%s\"", t)

	return []byte(result), nil
}
