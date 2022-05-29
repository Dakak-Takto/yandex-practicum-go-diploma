package entity

type Order struct {
	Number     int         `db:"number"      json:"number"`
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
