package middleware

import (
	"github.com/go-chi/chi"
)

func MiddlewareRoutes(c Middleware) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/getaccesstoken", c.GetAccessToken) //simulates the frontend returning an access token. should not be used in production
	return router
}
