package main

import (
	"auth_service/config"
	"auth_service/helper"
	routes "auth_service/route"
	"log"
	"net/http"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	router := routes.New()

	routes.RegisterRoutes(router)

	log.Println("🚀 Server running on :8080")

	for _, route := range router.Routes {
		log.Println(route.Method, route.Path)
	}

	log.Fatal(http.ListenAndServe(":8080", router.Router))
}
