package service

import "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"

type Service interface {
	RegisterUser(login string, password string) (*entity.User, error)
	AuthUser()
	CreateOrder()
	GetUserOrders()
	GetUserBalance()
	Withdraw()
	GetWithdrawals()
}
