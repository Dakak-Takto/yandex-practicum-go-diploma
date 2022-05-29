package entity

type Withdraw struct {
	Order       int          `db:"order_number" json:"order"`
	Sum         float64      `db:"sum" json:"sum"`
	UserID      int          `db:"user_id" json:"-"`
	ProcessedAt withdrawTime `db:"processed_at" json:"processed_at"`
}
