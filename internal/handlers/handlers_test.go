package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/require"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/mocks"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

func Test_handlers_userRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	m := mocks.NewMockService(ctrl)

	testUser := entity.User{
		ID:       24,
		Login:    "volodya",
		Password: utils.Hash("secret"),
		Balance:  2022,
	}

	m.EXPECT().GetUserByLogin(ctx, testUser.Login).Return(nil, nil)
	m.EXPECT().RegisterUser(ctx, testUser.Login, testUser.Password).Return(&testUser, nil)

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

	ctx := context.Background()

	m := mocks.NewMockService(ctrl)

	testUser := entity.User{
		ID:       24,
		Login:    "volodya",
		Password: utils.Hash("secret"),
		Balance:  2022,
	}

	m.EXPECT().AuthUser(ctx, testUser.Login, testUser.Password).Return(&testUser, nil)

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
