package storage

import (
	"context"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

//go:generate mockgen -destination=../mocks/mock_storage.go -package=mocks github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage Storage

type Storage interface {
	SaveUser(context.Context, *entity.User) (id int, err error)
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)

	SaveUserOrder(ctx context.Context, orderNumber string, userID int) (*entity.Order, error)
	SelectOrdersByUserID(ctx context.Context, userID int) ([]*entity.Order, error)
	GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) error
	UpdateUser(ctx context.Context, user *entity.User) error

	SaveWithdraw(context.Context, *entity.Withdraw) error
	SelectWithdrawals(ctx context.Context, userID int) ([]*entity.Withdraw, error)
	UserBalanceChange(ctx context.Context, userID int, delta float64) error

	SelectNewOrders(ctx context.Context) ([]*entity.Order, error)
}
