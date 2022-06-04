package handlers

import (
	"context"
	"net/http"
)

func (h *handler) CheckUserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessions.Get(r, cookieSessionName)
		if err != nil && !session.IsNew {
			log.Error("error get session: ", err)
			JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
			return
		}

		userID, exist := session.Values[cookieSessionUserIDKey].(int)

		if !exist {
			log.Error("пользователь не авторизован")
			JSONmsg(w, http.StatusUnauthorized, "error", "пользователь не авторизован")
			return
		}

		user, err := h.service.GetUserByID(r.Context(), userID)
		if err != nil {
			log.Errorf("пользователь %d не найден", userID)
			JSONmsg(w, http.StatusUnauthorized, "error", "пользователь не авторизован")
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) httpLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		recorder := &httpRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)

		e := r.Context().Value("error")
		if e != nil {
			log.Info(e)
		}

		log.Debugf("[%d] %s %s. ", recorder.Status, r.Method, r.RequestURI)
		log.Debugf("response: %s", recorder.response)
	})
}

type httpRecorder struct {
	http.ResponseWriter
	response []byte
	Status   int
}

func (h *httpRecorder) WriteHeader(status int) {
	h.Status = status
	h.ResponseWriter.WriteHeader(status)
}

func (h *httpRecorder) Write(b []byte) (int, error) {
	h.response = b
	return h.ResponseWriter.Write(b)
}
