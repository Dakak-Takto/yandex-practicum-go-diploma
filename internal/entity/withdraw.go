package entity

import (
	"fmt"
	"time"
)

type Withdraw struct {
	Order       string       `db:"order_number" json:"order"`
	Sum         float64      `db:"sum" json:"sum"`
	UserID      int          `db:"user_id" json:"-"`
	ProcessedAt withdrawTime `db:"processed_at" json:"processed_at"`
}

type withdrawTime time.Time

func (uploadedAt *withdrawTime) MarshalJSON() ([]byte, error) {

	t := time.Time(*uploadedAt).Format(time.RFC3339)
	result := fmt.Sprintf("\"%s\"", t)

	return []byte(result), nil
}
