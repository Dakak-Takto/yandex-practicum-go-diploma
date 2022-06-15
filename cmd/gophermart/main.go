package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	_accrual "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/client/accrual"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/config"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	_service "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
	_storage "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

func main() {
	log := logger.New()

	config.InitConfig()

	storage, err := _storage.NewPostgresStorage(config.DatabaseURI())
	if err != nil {
		log.Fatal(err)
	}
	service := _service.New(storage, log)

	cookieStore := initCookieStore(config.CookieStoreKey())

	handler := handlers.New(service, cookieStore, log)
	router := chi.NewRouter()
	handler.Register(router)

	accrual := _accrual.New(service, config.AccrualSystemAddress(), log)

	go accrual.Run(context.Background())

	log.Fatal(http.ListenAndServe(config.RunAddr(), router))

}

func initCookieStore(key string) *sessions.CookieStore {

	var keyPairs []byte
	var err error

	if len(key) == 0 {
		keyPairs, err = utils.Random(64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("new cookie secret key: %x", keyPairs)
	} else {
		keyPairs = []byte(key)
	}

	return sessions.NewCookieStore(keyPairs)
}
