package items

import (
	"github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func NewItemsRouter(s Service, m middleware.Middleware) *chi.Mux {
	router := chi.NewRouter()
	handler := newHandler(s)
	router.With(m.CheckMethod("POST"), m.Authorize(), m.AddHeader("Content-Type", "multipart/form-data")).Post("/item", handler.postItem)
	router.With(m.CheckMethod("PATCH"), m.Authorize()).Patch("/items/{id}", handler.updateItem)
	router.With(m.CheckMethod("DELETE"), m.Authorize()).Delete("/items/{id}", handler.deleteItem)
	router.With(m.CheckMethod("GET")).Get("/items/{id}", handler.getItem)
	router.With(m.CheckMethod("GET")).Get("/items", handler.getItems)
	return router
}
