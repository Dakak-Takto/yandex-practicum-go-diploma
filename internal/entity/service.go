package entity

type Service interface {
	CreateOrder(number int, userID int) (user *Order, err error)
	RegisterUser(login string, password string) (user *User, err error)
	AuthUser(login string, password string) (user *User, err error)
	GetUserByID(userID int) (user *User, err error)
	GetUserOrders(userID int) (orders []*Order, err error)
	GetUserByLogin(login string) (user *User, err error)
	UpdateUser(*User) error
	UpdateOrder(order *Order) error
	Withdraw(userID int, orderNumber int, sum float64) error
	GetWithdrawals(userID int) ([]*Withdraw, error)

	ProcessNewOrders() error
}
