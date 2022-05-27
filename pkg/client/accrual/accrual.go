package accrual

import "net/http"

// GET /api/orders/{number}` — получение информации о расчёте начислений баллов лояльности.
type Client interface {
	GetAccrualInfo(number int) (status string, accrual int, err error)
}

type client struct {
	httpClient *http.Client
}

func New() Client {
	return &client{
		httpClient: &http.Client{},
	}
}

func (s *client) GetAccrualInfo(orderNumber int) (status string, accrual int, err error) {
	return "", 0, nil
}
