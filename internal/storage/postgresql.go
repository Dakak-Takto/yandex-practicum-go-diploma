package storage

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

var _ entity.Storage = (*store)(nil)

type store struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func NewPostgresStorage(dsn string, logger *logrus.Logger) (entity.Storage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	store := store{
		db:  db,
		log: logger,
	}

	_, err = db.Exec(schema)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("store created")

	return &store, nil
}

func (s *store) SaveUser(user *entity.User) (*entity.User, error) {
	var userID int
	err := s.db.Get(&userID, `INSERT INTO users ( login, password ) VALUES ( $1, $2 ) RETURNING id`, user.Login, user.Password)

	if err != nil {
		s.log.Error("save user error:", err)
		return nil, err
	}

	return s.GetUserByID(userID)
}

func (s *store) SaveUserOrder(orderNumber string, userID int) (*entity.Order, error) {

	s.log.Debugf("insert order %s, user_id %d", orderNumber, userID)

	_, err := s.db.Exec(`INSERT INTO orders ( number, user_id ) VALUES ($1, $2)`, orderNumber, userID)
	if err != nil {
		s.log.Error("save order error:", err)
		return nil, err
	}

	return s.GetOrderByNumber(orderNumber)
}

func (s *store) SaveWithdraw(withdraw *entity.Withdraw) error {

	s.log.Debugf("insert withdraw: %+v", withdraw)

	_, err := s.db.NamedExec(`INSERT INTO withdrawals (order_number, sum, user_id) VALUES (:order_number, :sum, :user_id)`, withdraw)
	if err != nil {
		s.log.Error("save withdraw error: ", err)
		return err
	}

	return nil
}

func (s *store) GetOrderByNumber(number string) (*entity.Order, error) {

	var order entity.Order

	err := s.db.Get(&order, `SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE number = $1`, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		s.log.Error("get order by number error: ", err)
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
		s.log.Error("get user by login error: ", err)
		return nil, err
	}

	return &user, err
}

func (s *store) SelectOrdersByUserID(userID int) ([]*entity.Order, error) {

	var orders []*entity.Order

	err := s.db.Select(&orders, `SELECT number, accrual, status, user_id, uploaded_at FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		s.log.Error("select order by user_id error: ", err)
		return nil, err
	}

	return orders, nil
}

func (s *store) GetUserByID(id int) (*entity.User, error) {

	var user entity.User

	err := s.db.Get(&user, `SELECT id, login, password, balance FROM users WHERE users.id = $1`, id)
	if err != nil {
		s.log.Error("get user by id error: ", err)
		return nil, err
	}

	return &user, err
}

func (s *store) UpdateOrder(order *entity.Order) error {

	_, err := s.db.NamedExec(`UPDATE orders SET accrual=:accrual, status=:status, user_id=:user_id WHERE number=:number`, order)
	if err != nil {
		s.log.Error("update order error: ", err)
		return err
	}

	return nil
}

func (s *store) UpdateUser(user *entity.User) error {

	_, err := s.db.NamedExec(`UPDATE users SET login=:login, password=:password, balance=:balance WHERE id=:id`, user)
	if err != nil {
		s.log.Error("update user error: ", err)
		return err
	}

	return nil
}

func (s *store) SelectNewOrders() ([]*entity.Order, error) {

	var orders []*entity.Order

	statuses := []entity.OrderStatus{
		entity.OrderStatusNew,
		entity.OrderStatusRegistered,
		entity.OrderStatusProcessing,
	}

	err := s.db.Select(&orders, `SELECT number, accrual, status, user_id, uploaded_at FROM orders WHERE status = ANY($1)`, statuses)
	if err != nil {
		s.log.Error("select new orders error: ", err)
		return nil, err
	}
	return orders, nil
}

func (s *store) SelectWithdrawals(userID int) ([]*entity.Withdraw, error) {
	var withdrawals []*entity.Withdraw

	err := s.db.Select(&withdrawals, `SELECT order_number, sum, user_id, processed_at FROM withdrawals WHERE user_id = $1`, userID)
	if err != nil {
		s.log.Error("select withdrawals error: ", err)
		return nil, err
	}
	return withdrawals, nil
}

const schema string = `
CREATE TABLE IF NOT EXISTS users (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	login VARCHAR(128) NOT NULL, 
	password VARCHAR(64) NOT NULL,
	balance FLOAT NOT NULL DEFAULT (0)
);

CREATE TABLE IF NOT EXISTS orders (
	number VARCHAR(15) NOT NULL PRIMARY KEY, 
	status VARCHAR(20) NOT NULL DEFAULT ('NEW'),
	accrual FLOAT NOT NULL DEFAULT (0), 
	user_id INT NOT NULL REFERENCES users(id),
	uploaded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS withdrawals (
	order_number VARCHAR(15) NOT NULL PRIMARY KEY,
	sum FLOAT NOT NULL DEFAULT(0),
	user_id INT NOT NULL REFERENCES users(id),
	processed_at TIMESTAMP DEFAULT NOW()
);
`
