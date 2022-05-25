package service

import (
	"errors"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

type Service interface {
	RegisterUser(login string, password string) (user *entity.User, err error)
	AuthUser(login string, password string) (user *entity.User, err error)
	CreateOrder(number int, userID int) (user *entity.Order, err error)
	GetUserByID(userID int) (user *entity.User, err error)
	GetUserOrders(userID int) (orders []*entity.Order, err error)
	GetUserByLogin(login string) (user *entity.User, err error)
	Withdraw()
	GetWithdrawals()
}

var (
	ErrInvalidRequestFormat = errors.New("неверный формат запроса")
	ErrInvalidCredentials   = errors.New("неверная пара логин/пароль")
	ErrLoginAlreadyExists   = errors.New("логин уже занят")
	ErrInternalError        = errors.New("внутренняя ошибка")
	ErrCtxUserNotFound      = errors.New("пользователь не авторизован")
	ErrNotFound             = errors.New("error not found")
)
