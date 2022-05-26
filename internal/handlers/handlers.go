package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

type handler struct {
	service  entity.Service
	sessions *sessions.CookieStore
}

func New(service entity.Service, cookieStore *sessions.CookieStore) Handler {

	return &handler{
		service:  service,
		sessions: cookieStore,
	}
}

func (h *handler) Register(router chi.Router) {
	router.Route("/api/user/", func(r chi.Router) {

		r.MethodFunc(http.MethodPost, "/register", h.userRegister)
		r.MethodFunc(http.MethodPost, "/login", h.userLogin)

		r.Group(func(r chi.Router) {
			r.Use(h.CheckUserSession)

			r.MethodFunc(http.MethodGet, "/orders", h.userOrders)
			r.MethodFunc(http.MethodPost, "/orders", h.orderAdd)
			r.MethodFunc(http.MethodGet, "/balance", h.userBalance)
			r.MethodFunc(http.MethodPost, "/balance/withdraw", h.userBalanceWithdraw)
			r.MethodFunc(http.MethodGet, "/balance/withdrawals", h.userBalanceWithdrawals)
		})
	})
}

// регистрация пользователя
func (h *handler) userRegister(w http.ResponseWriter, r *http.Request) {

	var registerRequest UserRegisterDTO
	render.DecodeJSON(r.Body, &registerRequest)

	if registerRequest.Login == "" || registerRequest.Password == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "неверный формат запроса"})
		return
	}

	user, err := h.service.GetUserByLogin(registerRequest.Login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		log.Println(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
		return
	}

	if user != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, render.M{"error": "логин уже занят"})
		return
	}

	user, err = h.service.RegisterUser(registerRequest.Login, registerRequest.Password)
	if err != nil {
		log.Println("error register user:", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
		return
	}

	session, err := h.sessions.Get(r, cookieSessionName)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
	}
	session.Values[cookieSessionUserIDKey] = user.ID
	session.Save(r, w)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.M{"result": "пользователь зарегистрирован и аутентифицирован"})
	return

}

// аутентификация пользователя
func (h *handler) userLogin(w http.ResponseWriter, r *http.Request) {

	var loginRequest UserLoginDTO
	render.DecodeJSON(r.Body, &loginRequest)

	log.Println("логин запрос:", loginRequest)

	user, err := h.service.AuthUser(loginRequest.Login, loginRequest.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInternalError) {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
			return

		} else if errors.Is(err, entity.ErrInvalidRequestFormat) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"error": "неверный формат запроса"})
			return

		} else if errors.Is(err, entity.ErrInvalidCredentials) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, render.M{"error": "неверная пара логин/пароль"})
			return
		}
	}

	session, err := h.sessions.Get(r, cookieSessionName)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
	}
	session.Values[cookieSessionUserIDKey] = user.ID
	session.Save(r, w)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.M{"result": "пользователь успешно аутентифицирован"})
	return
}

// загрузка пользователем номера заказа для расчёта
func (h *handler) orderAdd(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "not auth", http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "interal server error", http.StatusInternalServerError)
	}

	orderNumber, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	_, err = h.service.CreateOrder(orderNumber, user.ID)
	if err != nil {
		http.Error(w, "error add order", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "order added")
}

//получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *handler) userOrders(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "not auth", http.StatusUnauthorized)
		return
	}

	orders, err := h.service.GetUserOrders(user.ID)
	if err != nil {
		log.Println("ошибка получения заказов из хранилища", err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, orders)
}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *handler) userBalance(w http.ResponseWriter, r *http.Request) {

}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *handler) userBalanceWithdraw(w http.ResponseWriter, r *http.Request) {

}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *handler) userBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {

}

func getUserFromContext(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value(userCtxKey("user")).(*entity.User)
	if !ok {
		return nil, entity.ErrCtxUserNotFound
	}
	return user, nil
}
