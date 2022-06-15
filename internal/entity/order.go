package entity

import (
	"encoding/json"
	"time"
)

type Order struct {
	Number     string      `db:"number" json:"number"`
	Accrual    float64     `db:"accrual" json:"accrual"`
	Status     OrderStatus `db:"status" json:"status"`
	UserID     int         `db:"user_id" json:"-"`
	UploadedAt time.Time   `db:"uploaded_at" json:"uploaded_at"`
}

//go:generate stringer -type=OrderStatus -trimprefix OrderStatus
type OrderStatus uint

const (
	OrderStatusUNKNOWN OrderStatus = iota
	OrderStatusNEW
	OrderStatusREGISTERED
	OrderStatusINVALID
	OrderStatusPROCESSING
	OrderStatusPROCESSED
)

func (o Order) MarshalJSON() ([]byte, error) {

	type OrderAlias Order

	aliasValue := struct {
		OrderAlias
		OrderStatus string `json:"status"`
		UploadedAt  string `json:"uploaded_at"`
	}{
		OrderAlias:  OrderAlias(o),
		UploadedAt:  o.UploadedAt.Format(time.RFC3339),
		OrderStatus: o.Status.String(),
	}

	return json.Marshal(aliasValue)
}
