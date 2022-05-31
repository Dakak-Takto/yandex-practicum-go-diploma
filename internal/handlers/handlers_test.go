package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/require"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/mocks"
)

func Test_handlers_userRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockService(ctrl)

	testUser := entity.User{
		ID:       24,
		Login:    "volodya",
		Password: utils.Hash("secret"),
		Balance:  2022,
	}

	m.EXPECT().GetUserByLogin(testUser.Login).Return(nil, nil)
	m.EXPECT().RegisterUser(testUser.Login, testUser.Password).Return(&testUser, nil)

	h := handler{
		service:  m,
		sessions: sessions.NewCookieStore([]byte("secret")),
	}
	handler := http.HandlerFunc(h.userRegister)

	b, _ := json.Marshal(map[string]string{
		"login":    testUser.Login,
		"password": testUser.Password,
	})
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(b))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	require.Equal(t, recorder.Code, http.StatusOK)
}

func Test_handlers_userLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockService(ctrl)

	testUser := entity.User{
		ID:       24,
		Login:    "volodya",
		Password: utils.Hash("secret"),
		Balance:  2022,
	}

	m.EXPECT().AuthUser(testUser.Login, testUser.Password).Return(&testUser, nil)

	h := handler{
		service:  m,
		sessions: sessions.NewCookieStore([]byte("secret")),
	}
	handler := http.HandlerFunc(h.userLogin)

	b, _ := json.Marshal(map[string]string{
		"login":    testUser.Login,
		"password": testUser.Password,
	})

	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(b))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	require.Equal(t, recorder.Code, http.StatusOK)
}
