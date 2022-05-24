package storage

type Storage interface {
	CreateUser()
	GetUserByLogin()
	GetOrderByUserID()
	GetUserBalance()
}
