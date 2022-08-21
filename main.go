package main

import (
	"fmt"

	"github.com/fnmzgdt/e_shop/src/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}
	router.StartServer()
}
