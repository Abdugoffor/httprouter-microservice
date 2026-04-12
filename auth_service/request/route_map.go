package request

type PermissionItem struct {
	Service string `json:"service"`
	Name    string `json:"name"`
}

type SyncPermissionsRequest struct {
	Permissions []PermissionItem `json:"permissions"`
}
