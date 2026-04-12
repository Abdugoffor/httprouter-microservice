package config

import (
	"context"
	"log"
)

func RunMigrations() {
	ctx := context.Background()

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS roles (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS permissions (
		id SERIAL PRIMARY KEY,
		service_name TEXT NOT NULL,
		name TEXT UNIQUE
	);

	CREATE TABLE IF NOT EXISTS user_roles (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		role_id INT REFERENCES roles(id) ON DELETE CASCADE,
		UNIQUE(user_id, role_id)
	);

	CREATE TABLE IF NOT EXISTS role_permissions (
		id SERIAL PRIMARY KEY,
		role_id INT REFERENCES roles(id) ON DELETE CASCADE,
		permission_id INT REFERENCES permissions(id) ON DELETE CASCADE,
		UNIQUE(role_id, permission_id)
	);
	`

	_, err := DB.Exec(ctx, query)
	if err != nil {
		log.Fatal("❌ Migration error:", err)
	}

	log.Println("✅ Migrations done")
}
