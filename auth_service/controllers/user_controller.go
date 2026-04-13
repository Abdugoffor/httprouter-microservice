package controllers

import (
	"auth_service/config"
	"auth_service/helper"
	"auth_service/request"
	response "auth_service/responce"
	"auth_service/services"
	"context"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(s services.AuthService) *AuthController {
	return &AuthController{service: s}
}

// REGISTER
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req request.RegisterRequest

	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, http.StatusOK, response.AuthResponse{Token: token})
}

// LOGIN
// role_id berilmasa -> roles ro'yxati qaytaradi
// role_id berilsa   -> token qaytaradi
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req request.LoginRequest

	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := c.service.Login(req.Email, req.Password, req.RoleID)
	if err != nil {
		JSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	JSON(w, http.StatusOK, result)
}

// SWITCH ROLE - mavjud token bilan boshqa rolega o'tish
func (c *AuthController) SwitchRole(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		JSON(w, http.StatusUnauthorized, "token required")
		return
	}

	tokenStr := strings.Replace(auth, "Bearer ", "", 1)
	claims, err := helper.ParseToken(tokenStr)
	if err != nil {
		JSON(w, http.StatusUnauthorized, "invalid token")
		return
	}

	var req request.SwitchRoleRequest
	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.service.SwitchRole(claims.UserID, req.RoleID)
	if err != nil {
		JSON(w, http.StatusForbidden, err.Error())
		return
	}

	JSON(w, http.StatusOK, response.AuthResponse{Token: token})
}

func (c *AuthController) SyncPermissions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req request.SyncPermissionsRequest

	if err := DecodeJSON(r, &req); err != nil {
		JSON(w, 400, err.Error())
		return
	}

	err := c.service.SyncPermissions(req)
	if err != nil {
		JSON(w, 500, err.Error())
		return
	}

	JSON(w, 200, "permissions synced")
}

func (c *AuthController) CheckPermission(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var req request.CheckPermissionRequest
	DecodeJSON(r, &req)

	permissionName := req.Method + " " + req.Path

	var exists bool

	err := config.DB.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1
			FROM role_permissions rp
			JOIN permissions p ON p.id = rp.permission_id
			WHERE rp.role_id = $1
			AND p.name = $2
		)
	`, req.RoleID, permissionName).Scan(&exists)

	if err != nil {
		JSON(w, 500, err.Error())
		return
	}

	JSON(w, 200, map[string]bool{
		"allowed": exists,
	})
}
