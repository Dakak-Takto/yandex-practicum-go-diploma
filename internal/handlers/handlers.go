package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
)

type handler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Register(router chi.Router) {
	router.Route("/api/user/", func(r chi.Router) {

		r.MethodFunc(http.MethodPost, "/register", h.registerUser)
		r.MethodFunc(http.MethodPost, "/login", h.loginUser)

		r.Group(func(r chi.Router) {
			r.MethodFunc(http.MethodGet, "/orders", nil)
			r.MethodFunc(http.MethodPost, "/orders", nil)
			r.MethodFunc(http.MethodGet, "/balance", nil)
			r.MethodFunc(http.MethodPost, "/balance/withdraw", nil)
			r.MethodFunc(http.MethodGet, "/balance/withdrawals", nil)
		})
	})
}

func (h *handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var registerRequest UserRegisterDTO
	render.DecodeJSON(r.Body, &registerRequest)
	user, err := h.service.RegisterUser(registerRequest.Login, registerRequest.Password)
	if err != nil {
		http.Error(w, "unable to register user", http.StatusInternalServerError)
	}
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, "user registered")
	log.Println(*user)
}

func (h *handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest UserLoginDTO
	render.DecodeJSON(r.Body, &loginRequest)
}
