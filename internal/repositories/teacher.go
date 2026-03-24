package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type TeacherRepository interface {
	Create(ctx context.Context, t *models.Teachers) error
	Update(ctx context.Context, t *models.Teachers) error
	GetByID(ctx context.Context, id int64) (*models.Teachers, error)
}

type teacherRepo struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) TeacherRepository {
	return &teacherRepo{
		db: db,
	}
}

func (r *teacherRepo) Create(ctx context.Context, t *models.Teachers) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *teacherRepo) Update(ctx context.Context, t *models.Teachers) error {
	// r.db.Model(&t).Select("IsActive").Updates(t).Error
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *teacherRepo) GetByID(ctx context.Context, id int64) (*models.Teachers, error) {
	var teacher models.Teachers

	err := r.db.WithContext(ctx).
		Preload("School").
		First(&teacher, id).Error

	if err != nil {
		return nil, err
	}

	return &teacher, nil
}
