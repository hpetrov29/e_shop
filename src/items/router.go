package items

import (
	M "github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func PostsRoutes(s Service, m M.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.With(m.CheckMethod("POST"), m.Authorize(), m.StaffAuthorize(), m.AddHeader("Content-Type", "application/json")).Post("/items", postItem(s))
	router.With(m.Authorize()).Post("/category", postCategory(s))
	router.With().Delete("/category", deleteCategory(s))
	router.With().Post("/brand", postBrand(s))
	router.With().Post("/size", postSizes(s))
	router.With().Delete("/size", deleteSizes(s))
	router.With().Post("/location", postLocations(s))
	router.With().Delete("/location", deleteLocations(s))
	router.With().Post("/discount", postDiscounts(s))
	router.With().Delete("/discount", deleteDiscounts(s))
	router.With().Post("/applydiscount", applyDiscounts(s))
	router.Patch("/items/{id}", updateItem(s))
	router.Delete("/items/{id}", deleteItem(s))
	router.Get("/items/{id}", getItem(s))
	router.Get("/items", getItems(s))
	return router
}
