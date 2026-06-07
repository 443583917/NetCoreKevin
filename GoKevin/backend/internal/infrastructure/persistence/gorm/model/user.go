package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type UserGORM struct {
	ID         int64                 `gorm:"primaryKey;autoIncrement"`
	UserName   string                `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password   string                `gorm:"type:varchar(100);not null"`
	RealName   string                `gorm:"type:varchar(50)"`
	Email      string                `gorm:"type:varchar(100)"`
	Phone      string                `gorm:"type:varchar(20)"`
	TenantID   int64                 `gorm:"index;not null"`
	IsDelete   soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime time.Time             `gorm:"autoCreateTime"`
	UpdateTime time.Time             `gorm:"autoUpdateTime"`
}

func (UserGORM) TableName() string {
	return "t_user"
}

type RoleGORM struct {
	ID         int64                 `gorm:"primaryKey;autoIncrement"`
	RoleName   string                `gorm:"type:varchar(50);not null"`
	RoleCode   string                `gorm:"type:varchar(50);uniqueIndex;not null"`
	TenantID   int64                 `gorm:"index;not null"`
	IsDelete   soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime time.Time             `gorm:"autoCreateTime"`
}

func (RoleGORM) TableName() string {
	return "t_role"
}

type PermissionGORM struct {
	ID       int64                 `gorm:"primaryKey;autoIncrement"`
	ParentID int64                 `gorm:"index;default:0"`
	Name     string                `gorm:"type:varchar(50);not null"`
	Code     string                `gorm:"type:varchar(100);uniqueIndex;not null"`
	Type     int                   `gorm:"not null"`
	URL      string                `gorm:"type:varchar(200)"`
	Icon     string                `gorm:"type:varchar(50)"`
	Sort     int                   `gorm:"default:0"`
	IsDelete soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
}

func (PermissionGORM) TableName() string {
	return "t_permission"
}

type UserRoleGORM struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey"`
}

func (UserRoleGORM) TableName() string {
	return "t_user_bind_role"
}

type RolePermissionGORM struct {
	RoleID       int64 `gorm:"primaryKey"`
	PermissionID int64 `gorm:"primaryKey"`
}

func (RolePermissionGORM) TableName() string {
	return "t_role_bind_permission"
}
