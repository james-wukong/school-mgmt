package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type AdminUserRepository interface {
	Create(ctx context.Context, t *models.AdminUser) error
	Update(ctx context.Context, t *models.AdminUser) error
	GetByID(ctx context.Context, id int64) (*models.AdminUser, error)
	GetSchoolByUserID(ctx context.Context, id int64) (*models.AdminUser, error)
}

type adminUserRepo struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepo{
		db: db,
	}
}

func (r *adminUserRepo) Create(ctx context.Context, t *models.AdminUser) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *adminUserRepo) Update(ctx context.Context, t *models.AdminUser) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *adminUserRepo) GetByID(ctx context.Context, id int64) (*models.AdminUser, error) {
	var u models.AdminUser

	err := r.db.
		WithContext(ctx).
		First(&u, id).Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *adminUserRepo) GetSchoolByUserID(
	ctx context.Context,
	id int64,
) (*models.AdminUser, error) {
	var u models.AdminUser

	err := r.db.
		WithContext(ctx).
		Preload("School").
		First(&u, id).
		Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}
