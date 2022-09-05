package middleware

import (
	"github.com/go-chi/chi"
)

func MiddlewareRoutes(c Controller) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/getaccesstoken", c.GetAccessToken)
	return router
}