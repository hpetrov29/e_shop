package users

import (
	"github.com/go-chi/chi"
)

func NewUsersRouter(s Service) *chi.Mux {
	router := chi.NewRouter()
	handler := newHandler(s)
	router.Post("/user", handler.registerUser)
	router.Post("/login", handler.login)
	return router
}
