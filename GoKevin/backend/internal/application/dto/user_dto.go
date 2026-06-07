package dto

type CreateUserRequest struct {
	UserName string `json:"userName" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"realName"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
}

type UserResponse struct {
	ID       int64          `json:"id"`
	UserName string         `json:"userName"`
	RealName string         `json:"realName"`
	Email    string         `json:"email"`
	Phone    string         `json:"phone"`
	TenantID int64          `json:"tenantId"`
	Roles    []RoleResponse `json:"roles,omitempty"`
}

type RoleResponse struct {
	ID       int64  `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
}

type UpdateUserRequest struct {
	RealName string `json:"realName"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
}

type ListUsersRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"pageSize,default=10" binding:"min=1,max=100"`
}
