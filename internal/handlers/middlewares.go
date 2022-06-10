package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

func (h *handler) CheckUserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessions.Get(r, cookieSessionName)
		if err != nil && !session.IsNew {
			h.log.Warn("error get session: ", err)
			JSONmsg(w, http.StatusInternalServerError, "error", "внутренняя ошибка сервера")
			return
		}

		userID, exist := session.Values[cookieSessionUserIDKey].(int)

		if !exist {
			h.log.Warn("пользователь не авторизован")
			JSONmsg(w, http.StatusUnauthorized, "error", "пользователь не авторизован")
			return
		}

		user, err := h.service.GetUserByID(r.Context(), userID)
		if err != nil {
			h.log.Warnf("пользователь %d не найден", userID)
			JSONmsg(w, http.StatusUnauthorized, "error", "пользователь не авторизован")
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) httpLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error(err)
			next.ServeHTTP(w, r)
		}
		r.Body.Close()

		reader := io.NopCloser(bytes.NewBuffer(body))
		r.Body = reader

		recorder := &httpRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)

		e := r.Context().Value("error")
		if e != nil {
			h.log.Info(e)
		}

		h.log.Debugf("[%d] %s %s %s", recorder.Status, r.Method, r.RequestURI, body)
		h.log.Debugf("response: %s", recorder.response)
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
