package migration

import (
	"log"

	"gorm.io/gorm"

	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence/gorm/model"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Run() error {
	log.Println("开始数据库迁移...")

	err := m.db.AutoMigrate(
		// User permission module
		&model.UserGORM{},
		&model.RoleGORM{},
		&model.PermissionGORM{},
		&model.UserRoleGORM{},
		&model.RolePermissionGORM{},

		// AI module
		&model.AIAppGORM{},
		&model.AIModelGORM{},
		&model.SkillGORM{},
		&model.AppSkillGORM{},
		&model.ChatSessionGORM{},
		&model.ChatMessageGORM{},
		&model.KnowledgeBaseGORM{},
		&model.DocumentGORM{},
	)

	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return err
	}

	log.Println("数据库迁移完成")
	return nil
}

func (m *Migrator) Seed() error {
	log.Println("开始初始化数据...")

	var count int64
	m.db.Model(&model.UserGORM{}).Where("user_name = ?", "admin").Count(&count)

	if count == 0 {
		admin := &model.UserGORM{
			UserName: "admin",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			RealName: "管理员",
			TenantID: 1000,
		}
		if err := m.db.Create(admin).Error; err != nil {
			return err
		}
		log.Println("默认管理员创建成功")
	}

	log.Println("初始化数据完成")
	return nil
}

func getModels() []string {
	return []string{
		"UserGORM",
		"RoleGORM",
		"PermissionGORM",
		"UserRoleGORM",
		"RolePermissionGORM",
		"AIAppGORM",
		"AIModelGORM",
		"SkillGORM",
		"AppSkillGORM",
		"ChatSessionGORM",
		"ChatMessageGORM",
		"KnowledgeBaseGORM",
		"DocumentGORM",
	}
}
