package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fnmzgdt/e_shop/src/middleware"
	"github.com/fnmzgdt/e_shop/src/posts"
	"github.com/fnmzgdt/e_shop/src/repositories"
	"github.com/fnmzgdt/e_shop/src/users"
	"github.com/fnmzgdt/e_shop/src/utils"
	"github.com/go-chi/chi"
)

func StartServer() *chi.Mux {
	var (
		port = utils.GetEnv("PORT", "8000")
		host = utils.GetEnv("DOCKER_HOST", "127.0.0.1")
	)

	mysql, err := repositories.SetupMySQLConnection()
	if err != nil {
		fmt.Println(err)
	}
	redis, err := repositories.SetupRedisConnection()
	if err != nil {
		fmt.Println(err)
	}

	postsService := posts.NewPostsService(mysql)
	usersService := users.NewUserssService(mysql, redis)
	middlewareController := middleware.NewMIddlewareController(redis)

	router := chi.NewRouter()

	router.Mount("/api/posts", posts.PostsRoutes(postsService, middlewareController))
	router.Mount("/api/users", users.UsersRoutes(usersService))
	router.Mount("/api/middleware", middleware.MiddlewareRoutes(middlewareController))

	fmt.Println("Server is listening on PORT " + port + ".")
	log.Fatal(http.ListenAndServe(host+":"+port, router))
	return router
}
