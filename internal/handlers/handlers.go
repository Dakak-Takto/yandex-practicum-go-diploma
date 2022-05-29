package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

type handler struct {
	service  entity.Service
	sessions *sessions.CookieStore
	log      *logrus.Logger
}

func New(service entity.Service, cookieStore *sessions.CookieStore, logger *logrus.Logger) Handler {

	return &handler{
		service:  service,
		sessions: cookieStore,
		log:      logger,
	}
}

func (h *handler) Register(router chi.Router) {
	router.Use(h.httpLog)

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

// POST /api/user/register. Регистрация пользователя
func (h *handler) userRegister(w http.ResponseWriter, r *http.Request) {

	// 200 — пользователь успешно зарегистрирован и аутентифицирован;
	// 400 — неверный формат запроса;
	// 409 — логин уже занят;
	// 500 — внутренняя ошибка сервера.

	var registerRequest UserRegisterDTO
	render.DecodeJSON(r.Body, &registerRequest)

	if registerRequest.Login == "" || registerRequest.Password == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "неверный формат запроса"})
		return
	}

	user, err := h.service.GetUserByLogin(registerRequest.Login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		h.log.Error(err)
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
		h.log.Error(err)
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
}

// POST /api/users/login. Aутентификация пользователя
func (h *handler) userLogin(w http.ResponseWriter, r *http.Request) {

	// 200 — пользователь успешно аутентифицирован;
	// 400 — неверный формат запроса;
	// 401 — неверная пара логин/пароль;
	// 500 — внутренняя ошибка сервера.

	var loginRequest UserLoginDTO
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

	user, err := getUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}

	orderNumber := string(body)

	order, err := h.service.GetOrderByNumber(orderNumber)
	if err != nil {
		if !errors.Is(err, entity.ErrNotFound) {
			h.log.Error(err)
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
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "новый номер заказа принят в обработку")
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
		h.log.Error(err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, orders)
}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *handler) userBalance(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "not auth", http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.service.GetWithdrawals(user.ID)
	if err != nil {
		h.log.Error("error get withdrawals", err)
		http.Error(w, "error get withdrawals", http.StatusInternalServerError)
		return
	}

	resp := struct {
		current   float64
		withdrawn float64
	}{}

	for _, w := range withdrawals {
		resp.withdrawn += w.Sum
	}
	resp.current = user.Balance

	render.JSON(w, r, resp)
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *handler) userBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r.Context())
	if err != nil {
		h.log.Error(err)
		http.Error(w, "not auth", http.StatusUnauthorized)
		return
	}

	var req WithdrawDTO
	render.DecodeJSON(r.Body, &req)
	err = h.service.Withdraw(user.ID, req.Order, req.Sum)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "error withdraw", http.StatusInternalServerError)
		return
	}
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *handler) userBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "not auth", http.StatusUnauthorized)
		return
	}
	withdrawals, err := h.service.GetWithdrawals(user.ID)
	if err != nil {
		h.log.Errorf("error get withdrawals: %s", err)
	}
	render.JSON(w, r, withdrawals)
}

func getUserFromContext(ctx context.Context) (*entity.User, error) {
	user, ok := ctx.Value(userCtxKey("user")).(*entity.User)
	if !ok {
		return nil, entity.ErrCtxUserNotFound
	}
	return user, nil
}
