package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/mocks"
)

func Test_service_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUserID := 12312

	ctx := context.Background()
	m := mocks.NewMockStorage(ctrl)

	m.EXPECT().GetUserByLogin(ctx, "login").
		Return(nil, nil)

	m.EXPECT().SaveUser(ctx, &entity.User{
		Login:    "login",
		Password: utils.Hash("password"),
	}).Return(newUserID, nil)

	m.EXPECT().GetUserByID(ctx, newUserID).
		Return(&entity.User{
			ID:       newUserID,
			Login:    "login",
			Password: utils.Hash("password"),
		}, nil)

	_, err := New(m).RegisterUser(ctx, "login", "password")
	require.NoError(t, err)
}

func Test_service_RegisterUser_LoginExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(ctx, "login").Return(&entity.User{}, nil)

	_, err := New(m).RegisterUser(ctx, "login", "password1")
	if assert.Error(t, err) {
		require.Equal(t, err, entity.ErrLoginAlreadyExists)
	}
}

func Test_service_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	const (
		login    string = "login"
		password string = "password"
	)

	u := entity.User{
		Login:    login,
		Password: utils.Hash(password),
	}

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(ctx, login).Return(&u, nil)

	_, err := New(m).AuthUser(ctx, login, password)
	require.NoError(t, err)
}

func Test_service_AuthUser_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	const (
		login    string = "login"
		password string = "password"
	)

	u := entity.User{
		Login:    login,
		Password: utils.Hash(password),
	}

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(ctx, login).Return(&u, nil)

	_, err := New(m).AuthUser(ctx, login, "wrong password")
	if assert.Error(t, err) {
		require.Equal(t, err, entity.ErrInvalidCredentials)
	}
}

func Test_service_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	m := mocks.NewMockStorage(ctrl)

	testOrder := entity.Order{
		Number: "422466257468622",
		UserID: 24,
	}
	m.EXPECT().GetOrderByNumber(ctx, testOrder.Number).Return(nil, nil)
	m.EXPECT().SaveUserOrder(ctx, testOrder.Number, testOrder.UserID).Return(&testOrder, nil)

	_, err := New(m).CreateOrder(ctx, testOrder.Number, testOrder.UserID)
	require.NoError(t, err)
}

func Test_service_CreateOrder_WrongNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	m := mocks.NewMockStorage(ctrl)

	testOrder := entity.Order{
		Number: "wrong number",
		UserID: 24,
	}

	_, err := New(m).CreateOrder(ctx, testOrder.Number, testOrder.UserID)
	if assert.Error(t, err) {
		require.Equal(t, err, entity.ErrOrderNumberIncorrect)
	}
}

func Test_service_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	m := mocks.NewMockStorage(ctrl)

	testOrders := []*entity.Order{
		{Number: "6174570601580", UserID: 24},
		{Number: "422466257468622", UserID: 24},
		{Number: "2352238521358", UserID: 24},
	}

	m.EXPECT().SelectOrdersByUserID(ctx, 24).Return(testOrders, nil)

	actualOrders, err := New(m).GetUserOrders(ctx, 24)
	require.NoError(t, err)
	require.Equal(t, testOrders, actualOrders)
}

func Test_service_GetWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	m := mocks.NewMockStorage(ctrl)

	userID := 24

	withdrawals := []*entity.Withdraw{
		{Order: "6174570601580", UserID: userID, Sum: 11},
		{Order: "422466257468622", UserID: userID, Sum: 22},
		{Order: "2352238521358", UserID: userID, Sum: 33},
	}

	m.EXPECT().SelectWithdrawals(ctx, userID).Return(withdrawals, nil)

	_, err := New(m).GetWithdrawals(ctx, userID)
	require.NoError(t, err)
}
