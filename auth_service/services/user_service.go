package services

import (
	"auth_service/config"
	"auth_service/helper"
	"auth_service/models"
	"auth_service/request"
	response "auth_service/responce"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(name, email, password string) (string, error)
	Login(email, password string, roleID *int) (interface{}, error)
	SwitchRole(userID int, roleID int) (string, error)
	SyncPermissions(req request.SyncPermissionsRequest) error
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

// REGISTER - user yaratadi va default 'user' roleini user_roles ga qo'shadi
func (s *authService) Register(name, email, password string) (string, error) {
	ctx := context.Background()

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	var user models.User

	err := config.DB.QueryRow(ctx,
		`INSERT INTO users (name, email, password)
		 VALUES ($1,$2,$3)
		 RETURNING id`,
		name, email, string(hashed),
	).Scan(&user.ID)

	if err != nil {
		return "", err
	}

	// default 'user' rolini topib user_roles ga qo'shamiz
	var defaultRoleID int
	err = config.DB.QueryRow(ctx,
		`SELECT id FROM roles WHERE name = 'user' LIMIT 1`,
	).Scan(&defaultRoleID)

	if err != nil {
		return "", errors.New("default role not found: " + err.Error())
	}

	_, err = config.DB.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		user.ID, defaultRoleID,
	)
	if err != nil {
		return "", err
	}

	return helper.GenerateToken(user.ID, defaultRoleID)
}

// LOGIN - role_id berilmasa roles ro'yxati, berilsa token qaytaradi
func (s *authService) Login(email, password string, roleID *int) (interface{}, error) {
	ctx := context.Background()

	var user models.User

	err := config.DB.QueryRow(ctx,
		`SELECT id, password FROM users WHERE email=$1`,
		email,
	).Scan(&user.ID, &user.Password)

	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	// foydalanuvchining barcha rollari
	rows, err := config.DB.Query(ctx,
		`SELECT r.id, r.name
		 FROM roles r
		 JOIN user_roles ur ON ur.role_id = r.id
		 WHERE ur.user_id = $1`,
		user.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []response.RoleItem
	for rows.Next() {
		var role response.RoleItem
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if len(roles) == 0 {
		return nil, errors.New("user has no roles assigned")
	}

	// role_id berilmagan - roles ro'yxatini qaytaramiz
	if roleID == nil {
		return response.LoginRolesResponse{
			UserID:  user.ID,
			Roles:   roles,
			Message: "select a role to get token",
		}, nil
	}

	// role_id berilgan - userni o'sha rolega egaligi tekshiriladi
	selectedRoleID := *roleID
	hasRole := false
	for _, r := range roles {
		if r.ID == selectedRoleID {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return nil, errors.New("user does not have this role")
	}

	token, err := helper.GenerateToken(user.ID, selectedRoleID)
	if err != nil {
		return nil, err
	}

	return response.AuthResponse{Token: token}, nil
}

// SWITCH ROLE - mavjud token bilan boshqa rolega o'tish
func (s *authService) SwitchRole(userID int, roleID int) (string, error) {
	ctx := context.Background()

	var exists bool
	err := config.DB.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM user_roles
			WHERE user_id = $1 AND role_id = $2
		)`,
		userID, roleID,
	).Scan(&exists)

	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("user does not have this role")
	}

	return helper.GenerateToken(userID, roleID)
}

func (s *authService) SyncPermissions(req request.SyncPermissionsRequest) error {
	ctx := context.Background()

	for _, p := range req.Permissions {

		_, err := config.DB.Exec(ctx,
			`INSERT INTO permissions(service_name, name)
			 VALUES ($1,$2)
			 ON CONFLICT(name) DO NOTHING`,
			p.Service, p.Name,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
