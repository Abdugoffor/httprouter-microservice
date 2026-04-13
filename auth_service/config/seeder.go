package config

import (
	"context"
	"log"
	"strings"
)

func RunRolePermissionSeeder() {
	ctx := context.Background()

	// 1. ROLES yaratamiz
	roles := []string{"admin", "mod", "user"}

	roleMap := map[string]int{}

	for _, role := range roles {
		var id int

		err := DB.QueryRow(ctx,
			`INSERT INTO roles (name)
			 VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id`, role).Scan(&id)

		if err != nil {
			log.Fatal("role error:", err)
		}

		roleMap[role] = id
	}

	// 2. DB dagi permissionslarni olamiz
	rows, err := DB.Query(ctx,
		`SELECT id, name FROM permissions`)
	if err != nil {
		log.Fatal("permission fetch error:", err)
	}
	defer rows.Close()

	type Perm struct {
		ID   int
		Name string
	}

	var permissions []Perm

	for rows.Next() {
		var p Perm
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			log.Fatal(err)
		}
		permissions = append(permissions, p)
	}

	// 3. ROLE → PERMISSION mapping 🔥

	for _, p := range permissions {

		method := strings.Split(p.Name, " ")[0] // GET /users → GET

		// ADMIN → ALL
		assign(roleMap["admin"], p.ID, ctx)

		// MOD → GET + POST
		if method == "GET" || method == "POST" {
			assign(roleMap["mod"], p.ID, ctx)
		}

		// USER → ONLY GET
		if method == "GET" {
			assign(roleMap["user"], p.ID, ctx)
		}
	}

	log.Println("✅ Roles + Permissions mapped successfully")
}

func assign(roleID, permID int, ctx context.Context) {
	_, err := DB.Exec(ctx,
		`INSERT INTO role_permissions (role_id, permission_id)
		 VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`,
		roleID, permID)

	if err != nil {
		log.Println("assign error:", err)
	}
}
