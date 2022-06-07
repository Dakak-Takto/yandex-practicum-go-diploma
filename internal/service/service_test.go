package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/mocks"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

func Test_service_RegisterUser(t *testing.T) {
	log := logger.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	m := mocks.NewMockStorage(ctrl)

	t.Run("normal register", func(t *testing.T) {

		newUserID := 12312

		m.EXPECT().
			GetUserByLogin(ctx, "login").
			Return(nil, nil)

		m.EXPECT().
			SaveUser(ctx, &entity.User{
				Login:    "login",
				Password: utils.Hash("password"),
			}).
			Return(newUserID, nil)

		m.EXPECT().
			GetUserByID(ctx, newUserID).
			Return(&entity.User{
				ID:       newUserID,
				Login:    "login",
				Password: utils.Hash("password"),
			}, nil)

		_, err := New(m, log).RegisterUser(ctx, "login", "password")
		require.NoError(t, err)
	})

	t.Run("login busy", func(t *testing.T) {

		m.EXPECT().
			GetUserByLogin(ctx, "login").
			Return(&entity.User{}, nil)

		_, err := New(m, log).RegisterUser(ctx, "login", "password1")
		if assert.Error(t, err) {
			require.Equal(t, err, entity.ErrLoginAlreadyExists)
		}
	})

	t.Run("auth", func(t *testing.T) {
		const (
			login    string = "login"
			password string = "password"
		)

		u := entity.User{
			Login:    login,
			Password: utils.Hash(password),
		}

		m.EXPECT().GetUserByLogin(ctx, login).Return(&u, nil)

		_, err := New(m, log).AuthUser(ctx, login, password)
		require.NoError(t, err)
	})

	t.Run("auth wrong credentials", func(t *testing.T) {
		const (
			login    string = "login"
			password string = "password"
		)

		u := entity.User{
			Login:    login,
			Password: utils.Hash(password),
		}

		m.EXPECT().GetUserByLogin(ctx, login).Return(&u, nil)

		_, err := New(m, log).AuthUser(ctx, login, "wrong password")
		if assert.Error(t, err) {
			require.Equal(t, err, entity.ErrInvalidCredentials)
		}
	})

	t.Run("create order", func(t *testing.T) {
		testOrder := entity.Order{
			Number: "422466257468622",
			UserID: 24,
		}
		m.EXPECT().GetOrderByNumber(ctx, testOrder.Number).Return(nil, nil)
		m.EXPECT().SaveUserOrder(ctx, testOrder.Number, testOrder.UserID).Return(&testOrder, nil)

		_, err := New(m, log).CreateOrder(ctx, testOrder.Number, testOrder.UserID)
		require.NoError(t, err)
	})

	t.Run("create order wrong number", func(t *testing.T) {
		testOrder := entity.Order{
			Number: "wrong number",
			UserID: 24,
		}

		_, err := New(m, log).CreateOrder(ctx, testOrder.Number, testOrder.UserID)
		if assert.Error(t, err) {
			require.Equal(t, err, entity.ErrOrderNumberIncorrect)
		}
	})

	t.Run("get user orders", func(t *testing.T) {
		testOrders := []*entity.Order{
			{Number: "6174570601580", UserID: 24},
			{Number: "422466257468622", UserID: 24},
			{Number: "2352238521358", UserID: 24},
		}

		m.EXPECT().
			SelectOrdersByUserID(ctx, 24).
			Return(testOrders, nil)

		actualOrders, err := New(m, log).GetUserOrders(ctx, 24)
		require.NoError(t, err)
		require.Equal(t, testOrders, actualOrders)
	})

	t.Run("get user withdrawals", func(t *testing.T) {
		userID := 24

		withdrawals := []*entity.Withdraw{
			{Order: "6174570601580", UserID: userID, Sum: 11},
			{Order: "422466257468622", UserID: userID, Sum: 22},
			{Order: "2352238521358", UserID: userID, Sum: 33},
		}

		m.EXPECT().
			SelectWithdrawals(ctx, userID).
			Return(withdrawals, nil)

		_, err := New(m, log).
			GetWithdrawals(ctx, userID)

		require.NoError(t, err)
	})

}
