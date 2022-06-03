package accrual

type (
	orderAccrualResponseDTO struct {
		Order   string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float64 `json:"accrual"`
	}
)
