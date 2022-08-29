package posts

import (
	"github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func PostsRoutes(s Service) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/post", middleware.Authorize(writePost(s)))
	return router
}
