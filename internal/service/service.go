package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/theplant/luhn"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

type service struct {
	storage storage.Storage
	log     *logrus.Logger
}

func New(storage storage.Storage, logger *logrus.Logger) Service {
	return &service{
		storage: storage,
		log:     logger,
	}
}

func (s *service) RegisterUser(ctx context.Context, login string, password string) (*entity.User, error) {

	if login == "" || password == "" {
		s.log.Warn("пустой логин или пароль")
		return nil, entity.ErrInvalidRequestFormat // если логин или пароль пустой
	}

	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		s.log.Warn("ошибка получения пользователя по логину")
		return nil, entity.ErrInternalError // ошибка при запросе их хранилища
	}
	if user != nil {
		s.log.Warn("логин занят")
		return nil, entity.ErrLoginAlreadyExists // логин занят
	}

	userID, err := s.storage.SaveUser(ctx, &entity.User{
		Login:    login,
		Password: utils.Hash(password),
	})

	if err != nil {
		s.log.Warn("ошибка сохранения пользователя", err)
		return nil, entity.ErrInternalError // ошибка при записи в хранилище
	}

	user, err = s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, entity.ErrInternalError
	}

	return user, nil
}

func (s *service) AuthUser(ctx context.Context, login string, password string) (*entity.User, error) {

	if login == "" || password == "" {
		s.log.Warn("логин или пароль пустой")
		return nil, entity.ErrInvalidRequestFormat // если логин или пароль пустой
	}

	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		s.log.Warnf("error get user: %s", err)
		return nil, err // ошибка при запросе из хранилища
	}

	if user.Password != utils.Hash(password) {
		s.log.Warn("invalid credentials")
		return nil, entity.ErrInvalidCredentials // пароль не совпадает или пользователя не существует
	}
	return user, nil
}

func (s *service) CreateOrder(ctx context.Context, number string, userID int) (*entity.Order, error) {
	orderInt, err := strconv.Atoi(number)
	if err != nil {
		return nil, entity.ErrOrderNumberIncorrect
	}
	if !luhn.Valid(orderInt) {
		return nil, entity.ErrOrderNumberIncorrect
	}

	order, err := s.GetOrderByNumber(ctx, number)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		s.log.Warn(err)
		return nil, entity.ErrInternalError
	}

	if order != nil {
		if order.UserID != userID {
			return nil, entity.ErrOrderNumberConflict
		}
		return nil, entity.ErrOrderNumberAlreadyExist
	}

	order, err = s.storage.SaveUserOrder(ctx, number, userID)
	if err != nil {
		s.log.Warnf("error save order: %s", err)
		return nil, err
	}
	return order, nil
}

func (s *service) UserBalanceChange(ctx context.Context, userID int, delta float64) error {

	if err := s.storage.UserBalanceChange(ctx, userID, delta); err != nil {
		return err
	}

	return nil
}

func (s *service) GetUserOrders(ctx context.Context, userID int) ([]*entity.Order, error) {
	orders, err := s.storage.SelectOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *service) Withdraw(ctx context.Context, userID int, orderNumber string, sum float64) error {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Warnf("Error get User: %s", err)
		return err
	}

	s.log.Debugf("user balance %f. withdraw %f", user.Balance, sum)

	err = s.storage.SaveWithdraw(ctx, &entity.Withdraw{
		UserID: user.ID,
		Sum:    sum,
		Order:  orderNumber,
	})
	if err != nil {
		s.log.Warnf("error save withdraw: %s", err)
	}

	err = s.UserBalanceChange(ctx, user.ID, -sum)

	return err
}

func (s *service) GetWithdrawals(ctx context.Context, userID int) ([]*entity.Withdraw, error) {
	withdrawals, err := s.storage.SelectWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (s *service) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.storage.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, user *entity.User) error {
	err := s.storage.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateOrder(ctx context.Context, order *entity.Order) error {
	err := s.storage.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error) {
	order, err := s.storage.GetOrderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) GetNewOrders(ctx context.Context) ([]*entity.Order, error) {
	orders, err := s.storage.SelectNewOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
