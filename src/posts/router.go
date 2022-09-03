package posts

import (
	M "github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func PostsRoutes(s Service, mc M.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/post", M.Authorize(writePost(s)))
	router.With(mc.Serialize).Post("/postt", writePost(s))
	return router
}
