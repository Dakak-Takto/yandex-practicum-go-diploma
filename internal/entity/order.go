package entity

type Order struct {
	Number     string      `db:"number"      json:"number"`
	Accrual    float64     `db:"accrual"     json:"accrual"`
	Status     OrderStatus `db:"status"      json:"status"`
	UserID     int         `db:"user_id"     json:"-"`
	UploadedAt orderTime   `db:"uploaded_at" json:"uploaded_at"`
}
