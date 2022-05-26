package entity

type Order struct {
	Number     int         `db:"number"      json:"number"`
	Accrual    int         `db:"accrual"     json:"accrual"`
	Status     orderStatus `db:"status"      json:"status"`
	UserID     int         `db:"user_id"     json:"-"`
	UploadedAt orderTime   `db:"uploaded_at" json:"uploaded_at"`
}
