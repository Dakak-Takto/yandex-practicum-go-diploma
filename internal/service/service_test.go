package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/stretchr/testify/require"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/mocks"
)

func Test_service_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)
	m.EXPECT().GetUserByLogin("login").Return(nil, nil)
	m.EXPECT().SaveUser(&entity.User{Login: "login", Password: utils.Hash("password")}).Return(&entity.User{Login: "login", Password: utils.Hash("password")}, nil)

	_, err := New(m).RegisterUser("login", "password")
	require.NoError(t, err)
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
	m.EXPECT().GetUserByLogin(login).Return(&u, nil)

	_, err := New(m).AuthUser(login, password)
	require.NoError(t, err)
}

func Test_service_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	testOrder := entity.Order{
		Number: "422466257468622",
		UserID: 24,
	}

	m.EXPECT().SaveUserOrder(testOrder.Number, testOrder.UserID).Return(&testOrder, nil)

	_, err := New(m).CreateOrder(testOrder.Number, testOrder.UserID)
	require.NoError(t, err)
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

	m.EXPECT().SelectOrdersByUserID(24).Return(testOrders, nil)

	actualOrders, err := New(m).GetUserOrders(24)
	require.NoError(t, err)
	require.Equal(t, testOrders, actualOrders)
}

func Test_service_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	testUser := &entity.User{
		ID:      24,
		Balance: 2022,
	}

	m.EXPECT().GetUserByID(testUser.ID).Return(testUser, nil)

	testWithdraw := &entity.Withdraw{
		UserID: testUser.ID,
		Sum:    22,
		Order:  "2352238521358",
	}

	wantBalance := testUser.Balance - testWithdraw.Sum

	m.EXPECT().SaveWithdraw(testWithdraw).Return(nil)
	m.EXPECT().UpdateUser(testUser).Return(nil)

	err := New(m).Withdraw(testUser.ID, testWithdraw.Order, testWithdraw.Sum)
	require.NoError(t, err)
	require.Equal(t, wantBalance, testUser.Balance)
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

	m.EXPECT().SelectWithdrawals(userID).Return(withdrawals, nil)

	_, err := New(m).GetWithdrawals(userID)
	require.NoError(t, err)
}
