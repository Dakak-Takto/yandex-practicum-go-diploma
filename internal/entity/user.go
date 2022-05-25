package entity

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Balance  int    `db:"balance"`
}

type UserService interface {
	Login(login string, password string) (*User, error)
	Register(login string, password string) (*User, error)
}

type UserStorage interface {
	GetByID(id int) (*User, error)
	GetByLogin(login string) (*User, error)
}
