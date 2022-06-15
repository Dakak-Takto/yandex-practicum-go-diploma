package entity

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	Order       string    `db:"order_number" json:"order"`
	Sum         float64   `db:"sum" json:"sum"`
	UserID      int       `db:"user_id" json:"-"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}

func (w Withdraw) MarshalJSON() ([]byte, error) {

	type WithdrawAlias Withdraw

	aliasValue := struct {
		WithdrawAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawAlias: WithdrawAlias(w),
		ProcessedAt:   w.ProcessedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}
