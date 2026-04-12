package response

type AuthResponse struct {
	Token string `json:"token"`
}

type RoleItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LoginRolesResponse struct {
	UserID  int        `json:"user_id"`
	Roles   []RoleItem `json:"roles"`
	Message string     `json:"message"`
}
