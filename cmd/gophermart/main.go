package main

import (
	"net/http"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	storage := storage.NewPostgresStorage("dsn string")
	service := service.New(storage)
	handler := handlers.New(service)

	router := chi.NewRouter()
	handler.Register(router)

	http.ListenAndServe("localhost:8080", router)
}
