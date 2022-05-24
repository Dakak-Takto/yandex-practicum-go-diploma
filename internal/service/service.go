package service

import (
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
)

type service struct {
	storage storage.Storage
}

func New(storage storage.Storage) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) RegisterUser(login string, password string) (*entity.User, error) {
	s.storage.CreateUser()
	return nil, nil
}

func (s *service) AuthUser() {
	return
}
func (s *service) CreateOrder() {
	return
}
func (s *service) GetUserOrders() {
	return
}
func (s *service) GetUserBalance() {
	return
}
func (s *service) Withdraw() {
	return
}
func (s *service) GetWithdrawals() {
	return
}
