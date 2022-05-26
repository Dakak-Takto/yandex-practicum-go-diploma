package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/pkg/client/accrual"
)

var _ entity.Service = (*service)(nil)

type service struct {
	storage       entity.Storage
	accrualClient accrual.Client
}

func New(storage entity.Storage) entity.Service {
	return &service{
		storage: storage,
	}
}

func (s *service) RegisterUser(login string, password string) (*entity.User, error) {

	if login == "" || password == "" {
		log.Println("пустой логин или пароль")
		return nil, entity.ErrInvalidRequestFormat // если логин или пароль пустой
	}

	user, err := s.storage.GetUserByLogin(login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		log.Println("ошибка получения пользователя по логину")
		return nil, entity.ErrInternalError // ошибка при запросе их хранилища
	}
	if user != nil {
		log.Println("логин занят")
		return nil, entity.ErrLoginAlreadyExists // логин занят
	}

	user, err = s.storage.SaveUser(&entity.User{
		Login:    login,
		Password: hashPassword(password),
	})

	if err != nil {
		log.Println("ошибка сохранения пользователя", err)
		return nil, entity.ErrInternalError // ошибка при записи в хранилище
	}

	return user, nil
}

func (s *service) AuthUser(login string, password string) (*entity.User, error) {

	if login == "" || password == "" {
		log.Println("логин или пароль пустой")
		return nil, entity.ErrInvalidRequestFormat // если логин или пароль пустой
	}

	user, err := s.storage.GetUserByLogin(login)
	if err != nil {
		return nil, entity.ErrInternalError // ошибка при запросе из хранилища
	}

	if user.Password != hashPassword(password) {
		return nil, entity.ErrInvalidCredentials // пароль не совпадает или пользователя не существует
	}
	return user, nil
}

func (s *service) CreateOrder(number int, userID int) (*entity.Order, error) {
	order, err := s.storage.SaveUserOrder(number, userID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *service) GetUserOrders(userID int) ([]*entity.Order, error) {
	return s.storage.SelectOrdersByUserID(userID)
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

func (s *service) GetUserByID(id int) (*entity.User, error) {
	return s.storage.GetUserByID(id)
}

func (s *service) GetUserByLogin(login string) (*entity.User, error) {
	return s.storage.GetUserByLogin(login)
}

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
