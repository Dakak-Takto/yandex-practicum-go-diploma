package entity

type User struct {
	ID       int    `db:"id"       json:"id"`
	Login    string `db:"login"    json:"login"`
	Password string `db:"password" json:"-"`
	Balance  int    `db:"balance"  json:"balance"`
}
