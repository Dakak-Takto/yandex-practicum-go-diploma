package entity

import (
	"fmt"
	"strings"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func (o *OrderStatus) UnmarshalJSON(b []byte) error {
	status := OrderStatus(strings.ToUpper(string(b)))

	switch status {
	case OrderStatusNew:
		b = []byte(OrderStatusNew)
		return nil

	case OrderStatusRegistered:
		b = []byte(OrderStatusRegistered)
		return nil

	case OrderStatusInvalid:
		b = []byte(OrderStatusInvalid)
		return nil

	case OrderStatusProcessed:
		b = []byte(OrderStatusProcessed)
		return nil

	case OrderStatusProcessing:
		b = []byte(OrderStatusProcessing)
		return nil
	}
	return fmt.Errorf("unknown order status")
}
