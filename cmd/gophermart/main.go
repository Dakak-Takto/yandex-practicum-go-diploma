package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	_accrual "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/client/accrual"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/config"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	_service "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

var log = logger.GetLoggerInstance()

func main() {
	config.InitConfig()

	store, err := storage.NewPostgresStorage(config.DatabaseURI())
	if err != nil {
		log.Fatal(err)
	}
	service := _service.New(store)

	cookieStore := initCookieStore(config.CookieStoreKey())

	handler := handlers.New(service, cookieStore)
	router := chi.NewRouter()
	handler.Register(router)

	accrual := _accrual.New(service, config.AccrualSystemAddress())

	go accrual.Run(context.Background())

	log.Fatal(http.ListenAndServe(config.RunAddr(), router))

}

func initCookieStore(key string) *sessions.CookieStore {

	var keyPairs []byte

	if len(key) == 0 {
		keyPairs = utils.Random(64)
		log.Infof("new cookie secret key: %x", keyPairs)
	} else {
		keyPairs = []byte(key)
	}

	return sessions.NewCookieStore(keyPairs)
}
