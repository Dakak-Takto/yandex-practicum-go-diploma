package storage

import "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"

type Storage interface {
	SaveUser(user *entity.User) (*entity.User, error)
	GetUserByLogin(login string) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)

	SaveUserOrder(orderNumber string, userID int) (*entity.Order, error)
	SelectOrdersByUserID(userID int) ([]*entity.Order, error)
	GetOrderByNumber(number string) (*entity.Order, error)
	UpdateOrder(order *entity.Order) error
	UpdateUser(user *entity.User) error

	SaveWithdraw(*entity.Withdraw) error
	SelectWithdrawals(userID int) ([]*entity.Withdraw, error)
	UserBalanceChange(userID int, delta float64) error

	SelectNewOrders() ([]*entity.Order, error)
}
