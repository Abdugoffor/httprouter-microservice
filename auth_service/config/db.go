package config

import (
	"auth_service/helper"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func DBConnect() *pgxpool.Pool {
	driver := helper.ENV("DB_DRIVER")
	if driver != "postgres" {
		log.Fatal("❌ pgx faqat PostgreSQL bilan ishlaydi")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		helper.ENV("DB_USER"),
		helper.ENV("DB_PASSWORD"),
		helper.ENV("DB_HOST"),
		helper.ENV("DB_PORT"), // ← 6432 (PgBouncer)
		helper.ENV("DB_NAME"),
		helper.ENV("DB_SSLMODE"),
		helper.ENV("DB_TIMEZONE"),
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	{
		if err != nil {
			log.Fatal("❌ DSN parse error:", err)
		}
	}

	cfg.MaxConns = 10 // PgBouncer allaqachon pool qiladi, shuning uchun past
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	{
		if err != nil {
			log.Fatal("❌ Failed to connect:", err)
		}
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("❌ DB ping error:", err)
	}

	log.Println("✅ Connected to PostgreSQL via PgBouncer 🚀")

	DB = db
	RunMigrations()
	RunRolePermissionSeeder()
	return db
}
