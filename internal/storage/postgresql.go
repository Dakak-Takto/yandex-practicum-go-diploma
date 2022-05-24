package storage

import (
	"database/sql"
)

type store struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) Storage {
	return &store{}
}

func (s *store) CreateUser() {
	return
}
func (s *store) GetUserByLogin() {
	return
}
func (s *store) GetOrderByUserID() {
	return
}
func (s *store) GetUserBalance() {
	return
}
