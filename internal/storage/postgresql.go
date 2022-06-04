package storage

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
)

var (
	_   Storage = (*store)(nil)
	log         = logger.GetLoggerInstance()
)

type store struct {
	db *sqlx.DB
}

func NewPostgresStorage(dsn string) (Storage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	store := store{
		db: db,
	}

	_, err = db.Exec(schema)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debug("connect database OK")

	return &store, nil
}

func (s *store) SaveUser(ctx context.Context, user *entity.User) (userID int, err error) {

	err = s.db.GetContext(ctx, &userID, `INSERT INTO users ( login, password ) VALUES ( $1, $2 ) RETURNING id`, user.Login, user.Password)

	if err != nil {
		log.Error("save user error:", err)
		return 0, err
	}

	return userID, err
}

func (s *store) SaveUserOrder(ctx context.Context, orderNumber string, userID int) (*entity.Order, error) {

	log.Debugf("insert order %s, user_id %d", orderNumber, userID)

	_, err := s.db.ExecContext(ctx, `INSERT INTO orders ( number, user_id ) VALUES ($1, $2)`, orderNumber, userID)
	if err != nil {
		log.Error("save order error:", err)
		return nil, err
	}

	return s.GetOrderByNumber(ctx, orderNumber)
}

func (s *store) SaveWithdraw(ctx context.Context, withdraw *entity.Withdraw) error {

	log.Debugf("insert withdraw: %+v", withdraw)

	_, err := s.db.NamedExecContext(ctx, `INSERT INTO withdrawals (order_number, sum, user_id) VALUES (:order_number, :sum, :user_id)`, withdraw)
	if err != nil {
		log.Error("save withdraw error: ", err)
		return err
	}

	return nil
}

func (s *store) GetOrderByNumber(ctx context.Context, number string) (*entity.Order, error) {

	var order entity.Order

	err := s.db.GetContext(ctx, &order, `SELECT number, status, accrual, user_id, uploaded_at FROM orders WHERE number = $1`, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		log.Error("get order by number error: ", err)
		return nil, err
	}

	return &order, err
}

func (s *store) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {

	var user entity.User

	err := s.db.GetContext(ctx, &user, `SELECT id, login, password, balance FROM users WHERE login = $1`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		log.Error("get user by login error: ", err)
		return nil, err
	}

	return &user, err
}

func (s *store) SelectOrdersByUserID(ctx context.Context, userID int) ([]*entity.Order, error) {

	var orders []*entity.Order

	err := s.db.SelectContext(ctx, &orders, `SELECT number, accrual, status, user_id, uploaded_at FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		log.Error("select order by user_id error: ", err)
		return nil, err
	}

	return orders, nil
}

func (s *store) GetUserByID(ctx context.Context, id int) (*entity.User, error) {

	var user entity.User

	err := s.db.GetContext(ctx, &user, `SELECT id, login, password, balance FROM users WHERE users.id = $1`, id)
	if err != nil {
		log.Error("get user by id error: ", err)
		return nil, err
	}

	return &user, err
}

func (s *store) UpdateOrder(ctx context.Context, order *entity.Order) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateOrder, err := s.db.PrepareNamedContext(ctx, `UPDATE orders SET accrual=:accrual, status=:status, user_id=:user_id WHERE number=:number`)
	if err != nil {
		log.Error("Prepare error:", err)
		return err
	}

	if _, err := tx.NamedStmtContext(ctx, updateOrder).ExecContext(ctx, order); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateUser(ctx context.Context, user *entity.User) error {

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	updateUser, err := s.db.PrepareNamedContext(ctx, `UPDATE users SET login=:login, password=:password, balance=:balance WHERE id=:id`)
	if err != nil {
		return err
	}

	if _, err := tx.NamedStmtContext(ctx, updateUser).Exec(user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *store) UserBalanceChange(ctx context.Context, userID int, delta float64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	updateUserBalance, err := s.db.PrepareContext(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`)
	if err != nil {
		return err
	}

	if _, err := tx.StmtContext(ctx, updateUserBalance).ExecContext(ctx, delta, userID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *store) SelectNewOrders(ctx context.Context) ([]*entity.Order, error) {

	var orders []*entity.Order

	statuses := []string{
		entity.OrderStatusNew,
		entity.OrderStatusRegistered,
		entity.OrderStatusProcessing,
	}

	err := s.db.SelectContext(ctx, &orders, `SELECT number, accrual, status, user_id, uploaded_at FROM orders WHERE status = ANY($1)`, statuses)
	if err != nil {
		log.Error("select new orders error: ", err)
		return nil, err
	}
	return orders, nil
}

func (s *store) SelectWithdrawals(ctx context.Context, userID int) ([]*entity.Withdraw, error) {
	var withdrawals []*entity.Withdraw

	err := s.db.SelectContext(ctx, &withdrawals, `SELECT order_number, sum, user_id, processed_at FROM withdrawals WHERE user_id = $1`, userID)
	if err != nil {
		log.Error("select withdrawals error: ", err)
		return nil, err
	}
	return withdrawals, nil
}
