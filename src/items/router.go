package items

import (
	M "github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func PostsRoutes(s Service, m M.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.With(m.CheckMethod("POST"), m.AddHeader("Content-Type", "multipart/form-data")).Post("/item", postItem(s))
	router.Patch("/items/{id}", updateItem(s))
	router.Delete("/items/{id}", deleteItem(s))
	router.Get("/items/{id}", getItem(s))
	router.Get("/items", getItems(s))
	return router
}
