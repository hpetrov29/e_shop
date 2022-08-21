package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fnmzgdt/e_shop/src/posts"
	"github.com/fnmzgdt/e_shop/src/repositories"
	"github.com/go-chi/chi"
)

func StartServer() *chi.Mux {
	port := os.Getenv("PORT")
	fmt.Println("Starting server on PORT " + port)

	db, err := repositories.SetupMySQLConnection()

	if err != nil {
		fmt.Println(err)
	}

	postsService := posts.NewPostsService(db)

	router := chi.NewRouter()
	router.Mount("/api/posts", posts.PostsRoutes(postsService))

	fmt.Println("Server is listening on PORT " + port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, router))
	return router
}
