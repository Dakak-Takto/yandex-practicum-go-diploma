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
