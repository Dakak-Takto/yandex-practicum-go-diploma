package service

import (
	"context"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

type Service interface {
	RegisterUser(ctx context.Context, login string, password string) (user *entity.User, err error)
	AuthUser(ctx context.Context, login string, password string) (user *entity.User, err error)
	GetUserByID(ctx context.Context, userID int) (user *entity.User, err error)
	GetUserByLogin(ctx context.Context, login string) (user *entity.User, err error)
	UpdateUser(context.Context, *entity.User) error

	CreateOrder(ctx context.Context, number string, userID int) (user *entity.Order, err error)
	GetUserOrders(ctx context.Context, userID int) (orders []*entity.Order, err error)
	UpdateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error)

	Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]*entity.Withdraw, error)
	UserBalanceChange(ctx context.Context, userID int, delta float64) error

	GetNewOrders(ctx context.Context) ([]*entity.Order, error)
}
