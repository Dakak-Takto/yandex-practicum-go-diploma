package entity

type Service interface {
	RegisterUser(login string, password string) (user *User, err error)
	AuthUser(login string, password string) (user *User, err error)
	GetUserByID(userID int) (user *User, err error)
	GetUserByLogin(login string) (user *User, err error)
	UpdateUser(*User) error

	CreateOrder(number string, userID int) (user *Order, err error)
	GetUserOrders(userID int) (orders []*Order, err error)
	UpdateOrder(order *Order) error
	GetOrderByNumber(number string) (*Order, error)

	Withdraw(userID int, orderNumber string, sum float64) error
	GetWithdrawals(userID int) ([]*Withdraw, error)
	ProcessNewOrders() error
}
