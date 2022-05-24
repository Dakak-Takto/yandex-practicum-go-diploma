package entity

import "context"

type Order struct {
	Number     int
	Accrual    int
	Status     string
	UserID     int
	UploadedAt int
}

type OrderService interface {
	Add(order *Order) (*Order, error)
}

type OrderStorage interface {
	Add(ctx context.Context, order *Order) (*Order, error)
}
