package entity

type User struct {
	ID       int
	Login    string
	Password string
	Balance  int
}

type UserService interface {
	Login(login string, password string) (*User, error)
	Register(login string, password string) (*User, error)
}

type UserStorage interface {
	GetByID(id int) (*User, error)
	GetByLogin(login string) (*User, error)
}
