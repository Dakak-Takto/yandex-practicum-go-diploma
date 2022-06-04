package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	if err := decodeJSON(r.Body, &registerRequest); err != nil {
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	if registerRequest.Login == "" || registerRequest.Password == "" {
		JSONmsg(w, http.StatusBadRequest, "error", "неверный формат запроса")
		return
	}

	user, err := h.service.GetUserByLogin(r.Context(), registerRequest.Login)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		log.Error(err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	if user != nil {
		log.Errorf("login %s already exists", user.Login)
		JSONmsg(w, http.StatusConflict, "error", "логин уже занят")
		return
	}

	user, err = h.service.RegisterUser(r.Context(), registerRequest.Login, registerRequest.Password)
	if err != nil {
		log.Error(err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	session, err := h.sessions.Get(r, cookieSessionName)
	if err != nil {
		log.Error(err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	session.Values[cookieSessionUserIDKey] = user.ID

	if err = session.Save(r, w); err != nil {
		log.Errorf("error save session: %s", err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	JSONmsg(w, http.StatusOK, "result", "пользователь зарегистрирован и аутентифицирован")
}

// POST /api/users/login. Аутентификация пользователя
func (h *handler) userLogin(w http.ResponseWriter, r *http.Request) {

	var loginRequest userLoginRequest
	if err := decodeJSON(r.Body, &loginRequest); err != nil {
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	user, err := h.service.AuthUser(r.Context(), loginRequest.Login, loginRequest.Password)
	if err != nil {

		if errors.Is(err, entity.ErrNotFound) {
			JSONmsg(w, http.StatusUnauthorized, "error", "неверная пара логин/пароль")
			return

		} else if errors.Is(err, entity.ErrInvalidRequestFormat) {
			JSONmsg(w, http.StatusBadRequest, "error", "неверный формат запроса")
			return

		}

		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	session, err := h.sessions.Get(r, cookieSessionName)
	if err != nil && !session.IsNew {
		log.Error(err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}
	session.Values[cookieSessionUserIDKey] = user.ID

	if err = session.Save(r, w); err != nil {
		log.Errorf("error save session: %s", err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	JSONmsg(w, http.StatusOK, "result", "пользователь успешно аутентифицирован")
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
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	orderNumber := string(body)

	if _, err = h.service.CreateOrder(r.Context(), orderNumber, user.ID); err != nil {

		if errors.Is(err, entity.ErrOrderNumberConflict) {
			JSONmsg(w, http.StatusConflict, "error", "номер заказа уже был загружен другим пользователем")
			return
		}

		if errors.Is(err, entity.ErrOrderNumberIncorrect) {
			JSONmsg(w, http.StatusUnprocessableEntity, "error", "неверный номер заказа")
			r.Context()
			return
		}

		if errors.Is(err, entity.ErrOrderNumberAlreadyExist) {
			JSONmsg(w, http.StatusOK, "result", "номер заказа уже был загружен")
			return
		}
	}
	JSONmsg(w, http.StatusAccepted, "result", "новый номер заказа принят в обработку")
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *handler) userOrders(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	orders, err := h.service.GetUserOrders(r.Context(), user.ID)
	if err != nil {
		log.Error(err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}
	JSON(w, http.StatusOK, orders)
}

// получение текущего баланса счёта баллов лояльности пользователя
func (h *handler) userBalance(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	withdrawals, err := h.service.GetWithdrawals(r.Context(), user.ID)
	if err != nil {
		log.Error("error get withdrawals", err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	var withdrawn float64 = 0

	for _, w := range withdrawals {
		withdrawn += w.Sum
	}

	JSON(w, http.StatusOK, balanceResponse{
		Current:   user.Balance,
		Withdrawn: withdrawn,
	})
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
func (h *handler) userBalanceWithdraw(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	var req withdrawRequest
	err := decodeJSON(r.Body, &req)
	if err != nil {
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}

	err = h.service.Withdraw(r.Context(), user.ID, req.Order, req.Sum)
	if err != nil {
		log.Error(err)
		JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
		return
	}
}

// получение информации о выводе средств с накопительного счёта пользователем
func (h *handler) userBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userCtxKey("user")).(*entity.User)

	withdrawals, err := h.service.GetWithdrawals(r.Context(), user.ID)
	if err != nil {
		log.Errorf("error get withdrawals: %s", err)
		return
	}

	JSON(w, http.StatusOK, withdrawals)
}
