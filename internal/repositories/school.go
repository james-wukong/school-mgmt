package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type SchoolRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Schools, error)
}

type schoolRepo struct {
	db *gorm.DB
}

func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolRepo{
		db: db,
	}
}

func (r *schoolRepo) GetByID(ctx context.Context, id int64) (*models.Schools, error) {
	var s models.Schools

	err := r.db.WithContext(ctx).
		First(&s, id).Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}
