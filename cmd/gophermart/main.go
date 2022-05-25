package main

import (
	"log"
	"net/http"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func main() {
	storage, err := storage.NewPostgresStorage("postgres://postgres:postgres@localhost/praktikum?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	service := service.New(storage)
	cookieStore := sessions.NewCookieStore([]byte("secret key"))
	handler := handlers.New(service, cookieStore)

	router := chi.NewRouter()
	handler.Register(router)

	http.ListenAndServe("localhost:8080", router)
}
