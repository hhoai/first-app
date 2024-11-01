package models

type PermissionsRequest struct {
	Permissions []int  `json:"permission" form:"permission"`
	RoleName    string `json:"role" form:"role"`
}

type PermissionWithRowNumber struct {
	RowNumber  int
	ID         int
	Permission string
}

type RolePermissionWithRowNumber struct {
	RowNumber      int
	ID             int
	RoleName       string
	PermissionName string
}
