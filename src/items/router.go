package items

import (
	"github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func PostsRoutes(s Service, m middleware.Middleware) *chi.Mux {
	router := chi.NewRouter()
	router.With(m.CheckMethod("POST"), m.Authorize(), m.AddHeader("Content-Type", "multipart/form-data")).Post("/item", postItem(s))
	router.With(m.CheckMethod("PATCH"), m.Authorize()).Patch("/items/{id}", updateItem(s))
	router.With(m.CheckMethod("DELETE"), m.Authorize()).Delete("/items/{id}", deleteItem(s))
	router.With(m.CheckMethod("GET")).Get("/items/{id}", getItem(s))
	router.With(m.CheckMethod("GET")).Get("/items", getItems(s))
	return router
}
