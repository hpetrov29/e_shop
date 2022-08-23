package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fnmzgdt/e_shop/src/posts"
	"github.com/fnmzgdt/e_shop/src/repositories"
	"github.com/fnmzgdt/e_shop/src/utils"
	"github.com/go-chi/chi"
)

func StartServer() *chi.Mux {
	var (
		port = utils.GetEnv("PORT", "8000")
		host = utils.GetEnv("GOLANG_ENV", "127.0.0.1")
	)

	db, err := repositories.SetupMySQLConnection()

	if err != nil {
		fmt.Println(err)
	}

	postsService := posts.NewPostsService(db)

	router := chi.NewRouter()
	router.Mount("/api/posts", posts.PostsRoutes(postsService))

	fmt.Println("Server is listening on PORT " + port + ".")
	log.Fatal(http.ListenAndServe(host+":"+port, router))
	return router
}
