package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
)

type handler struct {
	service  service.Service
	sessions *sessions.CookieStore
}

var log = logger.GetLoggerInstance()

func New(service service.Service, cookieStore *sessions.CookieStore) Handler {

	return &handler{
		service:  service,
		sessions: cookieStore,
	}
}

func (h *handler) Register(router chi.Router) {
	router.Use(h.httpLog)

	router.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", h.userRegister)
		r.Post("/login", h.userLogin)

		r.Group(func(r chi.Router) {
			r.Use(h.CheckUserSession)

			r.Get("/orders", h.userOrders)
			r.Post("/orders", h.orderAdd)
			r.Get("/balance", h.userBalance)
			r.Post("/balance/withdraw", h.userBalanceWithdraw)
			r.Get("/balance/withdrawals", h.userBalanceWithdrawals)
		})
	})
}

// POST /api/user/register. Регистрация пользователя
func (h *handler) userRegister(w http.ResponseWriter, r *http.Request) {

	var registerRequest userRegRequest
	render.DecodeJSON(r.Body, &registerRequest)

	if registerRequest.Login == "" || registerRequest.Password == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "неверный формат запроса"})
		return
	}

	user, err := h.service.GetUserByLogin(registerRequest.Login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		log.Error(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
		return
	}

	if user != nil {
		log.Errorf("login %s already exists", user.Login)
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, render.M{"error": "логин уже занят"})
		return
	}

	user, err = h.service.RegisterUser(registerRequest.Login, registerRequest.Password)
	if err != nil {
		log.Error(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
		return
	}

	session, err := h.sessions.Get(r, cookieSessionName)
	if err != nil {
		log.Error(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
		return
	}
	session.Values[cookieSessionUserIDKey] = user.ID
	session.Save(r, w)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.M{
		"result": "пользователь зарегистрирован и аутентифицирован",
	})
}

// POST /api/users/login. Aутентификация пользователя
func (h *handler) userLogin(w http.ResponseWriter, r *http.Request) {

	// 200 — пользователь успешно аутентифицирован;
	// 400 — неверный формат запроса;
	// 401 — неверная пара логин/пароль;
	// 500 — внутренняя ошибка сервера.

	var loginRequest userLoginRequest
	render.DecodeJSON(r.Body, &loginRequest)

	user, err := h.service.AuthUser(loginRequest.Login, loginRequest.Password)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, render.M{"error": "неверная пара логин/пароль"})
			return
		} else if errors.Is(err, entity.ErrInvalidRequestFormat) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"error": "неверный формат запроса"})
			return

		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
		return

	}

	session, err := h.sessions.Get(r, cookieSessionName)
	if err != nil {
		log.Error(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"error": "внутренняя ошибка сервера"})
	}
	session.Values[cookieSessionUserIDKey] = user.ID
	session.Save(r, w)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.M{"result": "пользователь успешно аутентифицирован"})
}

// POST /api/user/orders. Загрузка пользователем номера заказа для расчёта
func (h *handler) orderAdd(w http.ResponseWriter, r *http.Request) {

	// 200 — номер заказа уже был загружен этим пользователем;
	// 202 — новый номер заказа принят в обработку;
	// 400 — неверный формат запроса;
	// 401 — пользователь не аутентифицирован;
	// 409 — номер заказа уже был загружен другим пользователем;
	// 422 — неверный формат номера заказа;
	// 500 — внутренняя ошибка сервера.

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}

	orderNumber := string(body)

	order, err := h.service.GetOrderByNumber(orderNumber)
	if err != nil {

		if !errors.Is(err, entity.ErrNotFound) {
			log.Error(err)
			http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
	}

	if order != nil {

		if order.UserID != user.ID {
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, render.M{"error": "номер заказа уже был загружен другим пользователем"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, render.M{"error": "номер заказа уже был загружен"})

		return
	}

	_, err = h.service.CreateOrder(orderNumber, user.ID)
	if err != nil {
		if errors.Is(err, entity.ErrOrderNumberIncorrect) {
			render.Status(r, http.StatusUnprocessableEntity)
			render.PlainText(w, r, "неверный номер заказа")
		}
		render.Status(r, http.StatusInternalServerError)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "новый номер заказа принят в обработку")
}

//получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *handler) userOrders(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	orders, err := h.service.GetUserOrders(user.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, orders)
}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *handler) userBalance(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	withdrawals, err := h.service.GetWithdrawals(user.ID)
	if err != nil {
		log.Error("error get withdrawals", err)
		http.Error(w, "error get withdrawals", http.StatusInternalServerError)
		return
	}

	var withdrawn float64 = 0

	for _, w := range withdrawals {
		withdrawn += w.Sum
	}

	render.JSON(w, r, balanceResponse{
		Current:   user.Balance,
		Withdrawn: withdrawn,
	})
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *handler) userBalanceWithdraw(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	var req withdrawRequest
	render.DecodeJSON(r.Body, &req)
	err := h.service.Withdraw(user.ID, req.Order, req.Sum)
	if err != nil {
		log.Error(err)
		http.Error(w, "error withdraw", http.StatusInternalServerError)
		return
	}
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *handler) userBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	withdrawals, err := h.service.GetWithdrawals(user.ID)
	if err != nil {
		log.Errorf("error get withdrawals: %s", err)
	}
	render.JSON(w, r, withdrawals)
}
