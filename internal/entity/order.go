package entity

import (
	"fmt"
	"time"
)

type Order struct {
	Number     string      `db:"number"      json:"number"`
	Accrual    float64     `db:"accrual"     json:"accrual"`
	Status     OrderStatus `db:"status"      json:"status"`
	UserID     int         `db:"user_id"     json:"-"`
	UploadedAt orderTime   `db:"uploaded_at" json:"uploaded_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type orderTime time.Time

func (uploadedAt *orderTime) MarshalJSON() ([]byte, error) {

	t := time.Time(*uploadedAt).Format(time.RFC3339)
	result := fmt.Sprintf("\"%s\"", t)

	return []byte(result), nil
}
