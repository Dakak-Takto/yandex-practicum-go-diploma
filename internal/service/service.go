package service

import (
	"errors"
	"strconv"

	"github.com/theplant/luhn"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

var (
	_   entity.Service = (*service)(nil)
	log                = logger.GetLoggerInstance()
)

type service struct {
	storage entity.Storage
}

func New(storage entity.Storage) entity.Service {
	log.Debug("init service")
	return &service{
		storage: storage,
	}
}

func (s *service) RegisterUser(login string, password string) (*entity.User, error) {

	if login == "" || password == "" {
		log.Error("пустой логин или пароль")
		return nil, entity.ErrInvalidRequestFormat // если логин или пароль пустой
	}

	user, err := s.storage.GetUserByLogin(login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		log.Error("ошибка получения пользователя по логину")
		return nil, entity.ErrInternalError // ошибка при запросе их хранилища
	}
	if user != nil {
		log.Error("логин занят")
		return nil, entity.ErrLoginAlreadyExists // логин занят
	}

	user, err = s.storage.SaveUser(&entity.User{
		Login:    login,
		Password: utils.Hash(password),
	})

	if err != nil {
		log.Error("ошибка сохранения пользователя", err)
		return nil, entity.ErrInternalError // ошибка при записи в хранилище
	}

	return user, nil
}

func (s *service) AuthUser(login string, password string) (*entity.User, error) {

	if login == "" || password == "" {
		log.Error("логин или пароль пустой")
		return nil, entity.ErrInvalidRequestFormat // если логин или пароль пустой
	}

	user, err := s.storage.GetUserByLogin(login)
	if err != nil {
		log.Errorf("error get user: %s", err)
		return nil, err // ошибка при запросе из хранилища
	}

	if user.Password != utils.Hash(password) {
		log.Error("invalid credentials")
		return nil, entity.ErrInvalidCredentials // пароль не совпадает или пользователя не существует
	}
	return user, nil
}

func (s *service) CreateOrder(number string, userID int) (*entity.Order, error) {
	orderInt, err := strconv.Atoi(number)
	if err != nil {
		return nil, entity.ErrOrderNumberIncorrect
	}
	if !luhn.Valid(orderInt) {
		return nil, entity.ErrOrderNumberIncorrect
	}
	order, err := s.storage.SaveUserOrder(number, userID)
	if err != nil {
		log.Errorf("error save order: %s", err)
		return nil, err
	}
	return order, nil
}

func (s *service) GetUserOrders(userID int) ([]*entity.Order, error) {
	orders, err := s.storage.SelectOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *service) Withdraw(userID int, orderNumber string, sum float64) error {
	user, err := s.storage.GetUserByID(userID)
	if err != nil {
		log.Errorf("Error get User: %s", err)
		return err
	}

	log.Debugf("user balance %f. withdraw %f", user.Balance, sum)
	user.Balance = user.Balance - sum

	err = s.storage.SaveWithdraw(&entity.Withdraw{
		UserID: user.ID,
		Sum:    sum,
		Order:  orderNumber,
	})

	err = s.UpdateUser(user)

	return err
}

func (s *service) GetWithdrawals(userID int) ([]*entity.Withdraw, error) {
	withdrawals, err := s.storage.SelectWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (s *service) GetUserByID(id int) (*entity.User, error) {
	user, err := s.storage.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateUser(user *entity.User) error {
	err := s.storage.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetUserByLogin(login string) (*entity.User, error) {
	user, err := s.storage.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateOrder(order *entity.Order) error {
	err := s.storage.UpdateOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetOrderByNumber(number string) (*entity.Order, error) {
	order, err := s.storage.GetOrderByNumber(number)
	if err != nil {
		return nil, err
	}

	return order, nil
}
