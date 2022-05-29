package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (h *handler) CheckUserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := h.sessions.Get(r, cookieSessionName)
		if err != nil {
			fmt.Printf("error get session: %s", err)
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "")
			return
		}

		userID, exist := session.Values[cookieSessionUserIDKey].(int)

		if !exist {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, "userID not found")
			return
		}

		user, err := h.service.GetUserByID(userID)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.PlainText(w, r, fmt.Sprintf("user with id %d not found", userID))
			return
		}
		h.log.Debugf("user: %+v", *user)

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

		h.log.Debugf("[%d] %s %s. ", recorder.Status, r.Method, r.RequestURI)
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
