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

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(nil, "login").Return(nil, nil)
	m.EXPECT().SaveUser(nil, &entity.User{
		Login:    "login",
		Password: utils.Hash("password"),
	}).Return(&entity.User{Login: "login", Password: utils.Hash("password")}, nil)

	_, err := New(m).RegisterUser(nil, "login", "password")
	require.NoError(t, err)
}

func Test_service_RegisterUser_LoginExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(nil, "login").Return(&entity.User{}, nil)

	_, err := New(m).RegisterUser(nil, "login", "password1")
	if assert.Error(t, err) {
		require.Equal(t, err, entity.ErrLoginAlreadyExists)
	}
}

func Test_service_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		login    string = "login"
		password string = "password"
	)

	u := entity.User{
		Login:    login,
		Password: utils.Hash(password),
	}

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(nil, login).Return(&u, nil)

	_, err := New(m).AuthUser(nil, login, password)
	require.NoError(t, err)
}

func Test_service_AuthUser_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		login    string = "login"
		password string = "password"
	)

	u := entity.User{
		Login:    login,
		Password: utils.Hash(password),
	}

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin(nil, login).Return(&u, nil)

	_, err := New(m).AuthUser(nil, login, "wrong password")
	if assert.Error(t, err) {
		require.Equal(t, err, entity.ErrInvalidCredentials)
	}
}

func Test_service_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	testOrder := entity.Order{
		Number: "422466257468622",
		UserID: 24,
	}
	m.EXPECT().GetOrderByNumber(nil, testOrder.Number).Return(nil, nil)
	m.EXPECT().SaveUserOrder(nil, testOrder.Number, testOrder.UserID).Return(&testOrder, nil)

	_, err := New(m).CreateOrder(nil, testOrder.Number, testOrder.UserID)
	require.NoError(t, err)
}

func Test_service_CreateOrder_WrongNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	testOrder := entity.Order{
		Number: "wrong number",
		UserID: 24,
	}

	_, err := New(m).CreateOrder(nil, testOrder.Number, testOrder.UserID)
	if assert.Error(t, err) {
		require.Equal(t, err, entity.ErrOrderNumberIncorrect)
	}
}

func Test_service_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	testOrders := []*entity.Order{
		{Number: "6174570601580", UserID: 24},
		{Number: "422466257468622", UserID: 24},
		{Number: "2352238521358", UserID: 24},
	}

	m.EXPECT().SelectOrdersByUserID(nil, 24).Return(testOrders, nil)

	actualOrders, err := New(m).GetUserOrders(nil, 24)
	require.NoError(t, err)
	require.Equal(t, testOrders, actualOrders)
}

func Test_service_GetWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	userID := 24

	withdrawals := []*entity.Withdraw{
		{Order: "6174570601580", UserID: userID, Sum: 11},
		{Order: "422466257468622", UserID: userID, Sum: 22},
		{Order: "2352238521358", UserID: userID, Sum: 33},
	}

	m.EXPECT().SelectWithdrawals(gomock.Any(), userID).Return(withdrawals, nil)

	_, err := New(m).GetWithdrawals(context.Background(), userID)
	require.NoError(t, err)
}
