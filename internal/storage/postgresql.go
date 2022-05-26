package storage

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

var _ entity.Storage = (*store)(nil)

type store struct {
	db *sqlx.DB
}

func NewPostgresStorage(dsn string) (entity.Storage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	store := store{
		db: db,
	}

	//init tables here if need

	return &store, nil
}

func (s *store) SaveUser(user *entity.User) (*entity.User, error) {
	var userID int
	err := s.db.Get(&userID, `INSERT INTO users ( login, password ) VALUES ( $1, $2 ) RETURNING id`, user.Login, user.Password)

	if err != nil {
		return nil, err
	}
	return s.GetUserByID(userID)
}

func (s *store) SaveUserOrder(orderNumber int, userID int) (*entity.Order, error) {
	_, err := s.db.Exec(`INSERT INTO orders ( number, user_id ) VALUES ($1, $2)`, orderNumber, userID)
	if err != nil {
		return nil, err
	}
	return s.GetOrderByNumber(orderNumber)
}

func (s *store) GetOrderByNumber(number int) (*entity.Order, error) {
	var order entity.Order
	err := s.db.Get(&order, `SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE number = $1`, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &order, err
}

func (s *store) GetUserByLogin(login string) (*entity.User, error) {
	var user entity.User
	err := s.db.Get(&user, `SELECT id, login, password, balance FROM users WHERE login = $1`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &user, err
}

func (s *store) SelectOrdersByUserID(userID int) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := s.db.Select(&orders, `SELECT number, accrual, status, user_id, uploaded_at FROM orders WHERE user_id = $1`, userID)

	return orders, err
}

func (s *store) GetUserByID(id int) (*entity.User, error) {

	var user entity.User
	err := s.db.Get(&user, `SELECT id, login, password, balance FROM users WHERE users.id = $1`, id)

	return &user, err
}

const schema string = `

CREATE TABLE IF NOT EXISTS users (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	login VARCHAR(128) NOT NULL, 
	password VARCHAR(64) NOT NULL,
	balance INT NOT NULL DEFAULT (0)
);

CREATE TABLE IF NOT EXISTS orders (
	number INT NOT NULL PRIMARY KEY, 
	status SMALLINT NOT NULL DEFAULT (0),
	accrual INT NOT NULL DEFAULT (0), 
	user_id INT NOT NULL REFERENCES users(id),
	uploaded_at TIMESTAMP DEFAULT NOW()
);
`
