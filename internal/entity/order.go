package entity

import (
	"fmt"
	"time"
)

type Order struct {
	Number     string    `db:"number"      json:"number"`
	Accrual    float64   `db:"accrual"     json:"accrual"`
	Status     string    `db:"status"      json:"status"`
	UserID     int       `db:"user_id"     json:"-"`
	UploadedAt orderTime `db:"uploaded_at" json:"uploaded_at"`
}

const (
	OrderStatusNew        string = "NEW"
	OrderStatusRegistered string = "REGISTERED"
	OrderStatusInvalid    string = "INVALID"
	OrderStatusProcessing string = "PROCESSING"
	OrderStatusProcessed  string = "PROCESSED"
)

type orderTime time.Time

func (uploadedAt *orderTime) MarshalJSON() ([]byte, error) {

	t := time.Time(*uploadedAt).Format(time.RFC3339)
	result := fmt.Sprintf("\"%s\"", t)

	return []byte(result), nil
}
