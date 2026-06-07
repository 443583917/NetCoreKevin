package entity

import "time"

type User struct {
	ID         int64     `json:"id"`
	UserName   string    `json:"userName"`
	Password   string    `json:"-"`
	RealName   string    `json:"realName"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	TenantID   int64     `json:"tenantId"`
	IsDelete   bool      `json:"isDelete"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:t_user_bind_role;"`
}

type Role struct {
	ID         int64     `json:"id"`
	RoleName   string    `json:"roleName"`
	RoleCode   string    `json:"roleCode"`
	TenantID   int64     `json:"tenantId"`
	IsDelete   bool      `json:"isDelete"`
	CreateTime time.Time `json:"createTime"`

	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:t_role_bind_permission;"`
}

type Permission struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parentId"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Type     int    `json:"type"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	IsDelete bool   `json:"isDelete"`
}

type UserInfo struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"userId"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}
