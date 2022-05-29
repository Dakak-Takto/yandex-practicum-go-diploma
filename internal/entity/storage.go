package entity

type Storage interface {
	SaveUser(user *User) (*User, error)
	GetUserByLogin(login string) (*User, error)
	GetUserByID(id int) (*User, error)

	SaveUserOrder(orderNumber string, userID int) (*Order, error)
	SelectOrdersByUserID(userID int) ([]*Order, error)
	GetOrderByNumber(number string) (*Order, error)
	UpdateOrder(order *Order) error
	UpdateUser(user *User) error

	SaveWithdraw(*Withdraw) error
	SelectWithdrawals(userID int) ([]*Withdraw, error)

	SelectNewOrders() ([]*Order, error)
}
