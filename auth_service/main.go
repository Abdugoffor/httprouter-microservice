package main

import (
	"auth_service/clickhouse"
	"auth_service/config"
	"auth_service/helper"
	"auth_service/kafka"
	"auth_service/mq"
	routes "auth_service/route"
	"log"
	"net/http"
)

func main() {
	helper.LoadEnv()

	config.DBConnect()

	mq.Init()
	kafka.InitProducer()

	// ClickHouse table yaratish (agar yo'q bo'lsa)
	if err := clickhouse.InitAuditTable(); err != nil {
		log.Printf("⚠️  ClickHouse init warning: %v", err)

	} else {
		log.Println("✅ ClickHouse audit table ready")
	}

	router := routes.New()

	routes.RegisterRoutes(router)

	log.Println("🚀 Server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", router.Router))
}
