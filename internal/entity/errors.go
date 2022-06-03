package entity

import "errors"

var (
	ErrInvalidRequestFormat    = errors.New("неверный формат запроса")
	ErrInvalidCredentials      = errors.New("неверная пара логин/пароль")
	ErrLoginAlreadyExists      = errors.New("логин уже занят")
	ErrInternalError           = errors.New("внутренняя ошибка")
	ErrNotFound                = errors.New("not found")
	ErrOrderNumberIncorrect    = errors.New("неверный номер заказа")
	ErrOrderNumberConflict     = errors.New("заказ уже был загружен другим пользователем")
	ErrOrderNumberAlreadyExist = errors.New("заказ уже был загружен")
)
