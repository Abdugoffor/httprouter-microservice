package request

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   *int   `json:"role_id"` // optional: agar berilsa token qaytaradi, bermasa roles ro'yxati
}

type SwitchRoleRequest struct {
	RoleID int `json:"role_id"`
}

type CheckPermissionRequest struct {
	UserID int    `json:"user_id"`
	RoleID int    `json:"role_id"`
	Method string `json:"method"`
	Path   string `json:"path"`
}
