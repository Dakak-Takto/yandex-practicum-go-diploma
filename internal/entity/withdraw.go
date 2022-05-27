package entity

type Withdraw struct {
	Order       int          `db:"order_number" json:"order"`
	Sum         int          `db:"sum" json:"sum"`
	ProcessedAt withdrawTime `db:"processed_at" json:"processed_at"`
}
