package service

import "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"

type Service interface {
	RegisterUser(login string, password string) (user *entity.User, err error)
	AuthUser(login string, password string) (user *entity.User, err error)
	GetUserByID(userID int) (user *entity.User, err error)
	GetUserByLogin(login string) (user *entity.User, err error)
	UpdateUser(*entity.User) error

	CreateOrder(number string, userID int) (user *entity.Order, err error)
	GetUserOrders(userID int) (orders []*entity.Order, err error)
	UpdateOrder(order *entity.Order) error
	GetOrderByNumber(number string) (*entity.Order, error)

	Withdraw(userID int, orderNumber string, sum float64) error
	GetWithdrawals(userID int) ([]*entity.Withdraw, error)
	ProcessNewOrders() error
}
