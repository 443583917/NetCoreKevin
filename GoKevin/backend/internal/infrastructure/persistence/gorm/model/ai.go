package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type AIAppGORM struct {
	ID           int64                 `gorm:"primaryKey;autoIncrement"`
	AppName      string                `gorm:"type:varchar(100);not null"`
	AppDesc      string                `gorm:"type:varchar(500)"`
	ModelID      int64                 `gorm:"index"`
	SystemPrompt string                `gorm:"type:text"`
	TenantID     int64                 `gorm:"index;not null"`
	IsDelete     soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime   time.Time             `gorm:"autoCreateTime"`
}

func (AIAppGORM) TableName() string {
	return "t_ai_apps"
}

type AIModelGORM struct {
	ID         int64                 `gorm:"primaryKey;autoIncrement"`
	ModelName  string                `gorm:"type:varchar(100);not null"`
	Provider   string                `gorm:"type:varchar(50);not null"`
	APIKey     string                `gorm:"type:varchar(200)"`
	BaseURL    string                `gorm:"type:varchar(200)"`
	MaxTokens  int                   `gorm:"default:4096"`
	IsDelete   soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime time.Time             `gorm:"autoCreateTime"`
}

func (AIModelGORM) TableName() string {
	return "t_ai_models"
}

type SkillGORM struct {
	ID          int64                 `gorm:"primaryKey;autoIncrement"`
	SkillName   string                `gorm:"type:varchar(100);not null"`
	SkillCode   string                `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string                `gorm:"type:varchar(500)"`
	SkillType   string                `gorm:"type:varchar(20);not null"`
	Config      string                `gorm:"type:json"`
	IsDelete    soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime  time.Time             `gorm:"autoCreateTime"`
}

func (SkillGORM) TableName() string {
	return "t_ai_skills"
}

type AppSkillGORM struct {
	AppID   int64 `gorm:"primaryKey"`
	SkillID int64 `gorm:"primaryKey"`
}

func (AppSkillGORM) TableName() string {
	return "t_ai_apps_bind_skill"
}

type ChatSessionGORM struct {
	ID         int64                 `gorm:"primaryKey;autoIncrement"`
	UserID     int64                 `gorm:"index;not null"`
	AppID      int64                 `gorm:"index;not null"`
	Title      string                `gorm:"type:varchar(200)"`
	IsDelete   soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime time.Time             `gorm:"autoCreateTime"`
	UpdateTime time.Time             `gorm:"autoUpdateTime"`
}

func (ChatSessionGORM) TableName() string {
	return "t_ai_chat_sessions"
}

type ChatMessageGORM struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	SessionID  int64     `gorm:"index;not null"`
	Role       string    `gorm:"type:varchar(20);not null"`
	Content    string    `gorm:"type:text;not null"`
	Tokens     int       `gorm:"default:0"`
	CreateTime time.Time `gorm:"autoCreateTime"`
}

func (ChatMessageGORM) TableName() string {
	return "t_ai_chat_messages"
}

type KnowledgeBaseGORM struct {
	ID          int64                 `gorm:"primaryKey;autoIncrement"`
	Name        string                `gorm:"type:varchar(100);not null"`
	Description string                `gorm:"type:varchar(500)"`
	VectorModel string                `gorm:"type:varchar(50)"`
	TenantID    int64                 `gorm:"index;not null"`
	IsDelete    soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime  time.Time             `gorm:"autoCreateTime"`
}

func (KnowledgeBaseGORM) TableName() string {
	return "t_ai_knowledge_bases"
}

type DocumentGORM struct {
	ID         int64                 `gorm:"primaryKey;autoIncrement"`
	KBID       int64                 `gorm:"index;not null"`
	FileName   string                `gorm:"type:varchar(200);not null"`
	FileURL    string                `gorm:"type:varchar(500)"`
	Status     string                `gorm:"type:varchar(20);default:pending"`
	ChunkCount int                   `gorm:"default:0"`
	IsDelete   soft_delete.DeletedAt `gorm:"softDelete:flag;index"`
	CreateTime time.Time             `gorm:"autoCreateTime"`
}

func (DocumentGORM) TableName() string {
	return "t_ai_documents"
}
