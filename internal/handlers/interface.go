package handlers

import "github.com/go-chi/chi/v5"

type Handler interface {
	Register(router chi.Router)
}

type userCtxKey string

const (
	cookieSessionName      string = "token"
	cookieSessionUserIDKey string = "user_id"
)
