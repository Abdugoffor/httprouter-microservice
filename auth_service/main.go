package main

import (
	"auth_service/clickhouse"
	"auth_service/config"
	"auth_service/helper"
	routes "auth_service/route"
	"log"
	"net/http"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	// ClickHouse table yaratish (agar yo'q bo'lsa)
	if err := clickhouse.InitAuditTable(); err != nil {
		log.Printf("⚠️  ClickHouse init warning: %v", err)
		// Fatal qilmaymiz — audit log ishlamasa ham asosiy servis ishlayversin
	} else {
		log.Println("✅ ClickHouse audit table ready")
	}

	router := routes.New()

	routes.RegisterRoutes(router)

	log.Println("🚀 Server running on :8080")

	// for _, route := range router.Routes {
	// 	log.Println(route.Method, route.Path)
	// }

	log.Fatal(http.ListenAndServe(":8080", router.Router))
}
