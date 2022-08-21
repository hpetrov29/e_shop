package posts

import "github.com/go-chi/chi"

func PostsRoutes(s Service) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/post", writePost(s))
	return router
}