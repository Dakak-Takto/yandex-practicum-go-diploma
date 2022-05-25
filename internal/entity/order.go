package entity

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Order struct {
	Number     int         `db:"number" json:"number"`
	Accrual    int         `db:"accrual" json:"accrual"`
	Status     orderStatus `db:"status" json:"status"`
	UserID     int         `db:"user_id" json:"-"`
	UploadedAt orderTime   `db:"uploaded_at" json:"uploaded_at"`
}

type OrderService interface {
	Add(order *Order) (*Order, error)
}

type OrderStorage interface {
	Add(ctx context.Context, order *Order) (*Order, error)
}

type orderTime time.Time

func (uploadedAt *orderTime) MarshalJSON() ([]byte, error) {
	result := fmt.Sprintf("\"%s\"", time.Time(*uploadedAt).Format(time.RFC3339))
	return []byte(result), nil
}

type orderStatus uint8

const (
	orderStatusNew orderStatus = 0 + iota
	orderStatusProcessing
	orderStatusInvalid
	orderStatusProcessed
)

func (orderStatus *orderStatus) MarshalJSON() ([]byte, error) {
	var result string

	switch *orderStatus {
	case orderStatusNew:
		result = fmt.Sprintf("\"%s\"", "new")
	case orderStatusProcessing:
		result = fmt.Sprintf("\"%s\"", "processing")
	case orderStatusInvalid:
		result = fmt.Sprintf("\"%s\"", "invalid")
	case orderStatusProcessed:
		result = fmt.Sprintf("\"%s\"", "processed")
	default:
		return nil, errors.New("unknown order status")
	}
	return []byte(result), nil
}
