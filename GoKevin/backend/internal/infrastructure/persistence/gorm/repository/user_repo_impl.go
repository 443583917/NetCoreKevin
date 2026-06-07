package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kevin-ai/go-kevin/internal/domain/entity"
	"github.com/kevin-ai/go-kevin/internal/domain/repository"
	"github.com/kevin-ai/go-kevin/internal/infrastructure/persistence/gorm/model"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	gormModel := r.toGORMModel(user)
	return r.db.WithContext(ctx).Create(gormModel).Error
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var m model.UserGORM
	err := r.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	gormModel := r.toGORMModel(user)
	return r.db.WithContext(ctx).Save(gormModel).Error
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.UserGORM{}, id).Error
}

func (r *UserRepositoryImpl) GetByUserName(ctx context.Context, userName string) (*entity.User, error) {
	var m model.UserGORM
	err := r.db.WithContext(ctx).Where("user_name = ?", userName).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserGORM
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *UserRepositoryImpl) ListByTenantID(ctx context.Context, tenantID int64, page, pageSize int) ([]*entity.User, int64, error) {
	var models []model.UserGORM
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	query.Model(&model.UserGORM{}).Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(models))
	for i, m := range models {
		users[i] = r.toEntity(&m)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) GetWithRoles(ctx context.Context, id int64) (*entity.User, error) {
	var m model.UserGORM
	err := r.db.WithContext(ctx).
		Preload("Roles").
		First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&m), nil
}

func (r *UserRepositoryImpl) BindRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.UserRoleGORM{}).Error
	if err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		ur := &model.UserRoleGORM{
			UserID: userID,
			RoleID: roleID,
		}
		if err := r.db.WithContext(ctx).Create(ur).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepositoryImpl) toGORMModel(user *entity.User) *model.UserGORM {
	return &model.UserGORM{
		ID:       user.ID,
		UserName: user.UserName,
		Password: user.Password,
		RealName: user.RealName,
		Email:    user.Email,
		Phone:    user.Phone,
		TenantID: user.TenantID,
	}
}

func (r *UserRepositoryImpl) toEntity(m *model.UserGORM) *entity.User {
	return &entity.User{
		ID:         m.ID,
		UserName:   m.UserName,
		RealName:   m.RealName,
		Email:      m.Email,
		Phone:      m.Phone,
		TenantID:   m.TenantID,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}
