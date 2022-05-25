package storage

import (
	"errors"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

type Storage interface {
	SaveUser(user *entity.User) (*entity.User, error)
	SaveUserOrder(orderNumber int, userID int) (*entity.Order, error)
	GetUserByLogin(login string) (*entity.User, error)
	SelectOrdersByUserID(userID int) ([]*entity.Order, error)
	GetUserByID(id int) (*entity.User, error)
	GetOrderByNumber(number int) (*entity.Order, error)
}

var (
	ErrNotFound = errors.New("not found")
)
