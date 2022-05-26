package accrual

import "net/http"

// GET /api/orders/{number}` — получение информации о расчёте начислений баллов лояльности.
type Client interface {
	AddOrders(orders ...int)
	Run() error
}

type client struct {
	httpClient *http.Client
}

func New() Client {
	return &client{
		httpClient: &http.Client{},
	}
}

func (c *client) AddOrders(orders ...int) {

}

func (c *client) Run() error {
	return nil
}
