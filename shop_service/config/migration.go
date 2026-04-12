package config

import (
	"context"
	"log"
)

func RunMigrations() {
	ctx := context.Background()

	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT,
		price NUMERIC,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := DB.Exec(ctx, query)
	if err != nil {
		log.Fatal("❌ Migration error:", err)
	}

	log.Println("✅ Migrations done")
}
