package users

import (
	"github.com/go-chi/chi"
)

func UsersRoutes(s Service) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/user", registerUser(s))
	return router
}
