package main

import (
	"log"
	"net/http"
	"shop_service/config"
	"shop_service/helper"
	"shop_service/kafka"
	"shop_service/mq"
	routes "shop_service/route"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	go mq.Consume()

	go kafka.Consume()

	router := routes.New()

	routes.RegisterRoutes(router)

	log.Println("🚀 Server running on :8081")

	// for _, route := range router.Routes {
	// 	log.Println(route.Method, route.Path)
	// }
	router.SyncRoutes()

	log.Fatal(http.ListenAndServe(":8081", router.Router))
}
