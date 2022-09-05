package posts

import (
	M "github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/go-chi/chi"
)

func PostsRoutes(s Service, mc M.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.With(mc.Serialize, mc.Authorize).Post("/post", writePost(s))
	return router
}
