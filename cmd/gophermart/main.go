package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/config"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/utils"
)

var log = logger.GetLoggerInstance()

func main() {

	storage, err := storage.NewPostgresStorage(config.DatabaseURI())
	if err != nil {
		log.Fatal(err)
	}
	service := service.New(storage)

	cookieStore := initCookieStore(config.CookieStoreKey())

	handler := handlers.New(service, cookieStore)
	router := chi.NewRouter()
	handler.Register(router)

	log.Infof("lister %s", config.RunAddr())
	go http.ListenAndServe(config.RunAddr(), router)

	for {
		time.Sleep(time.Second)
		err := service.ProcessNewOrders()
		if err != nil {
			log.Error(err)
		}
	}
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
