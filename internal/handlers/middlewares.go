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

		ctx := context.WithValue(r.Context(), userCtxKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
