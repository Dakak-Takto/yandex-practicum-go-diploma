package main

import (
	"fmt"
	"net/http"
	"path"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
)

func main() {
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:05:05",
		FullTimestamp:   true,
		ForceColors:     true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {

			return "", fmt.Sprintf(" %s:%d", path.Base(f.File), f.Line)
		},
	})
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)

	log.Info("init storage")
	storage, err := storage.NewPostgresStorage("postgres://postgres:postgres@localhost/praktikum?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	log.Info("init service")
	service := service.New(storage, log)

	log.Info("init cookiestore")
	cookieStore := sessions.NewCookieStore([]byte("secret key"))

	log.Info("init handler")
	handler := handlers.New(service, cookieStore, log)

	log.Info("init router")
	router := chi.NewRouter()

	log.Info("register handler")
	handler.Register(router)

	log.Info("lister and serve http")
	http.ListenAndServe("localhost:8080", router)
}
