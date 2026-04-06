// Package repositories
package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type TeacherSubjectRepository interface {
	Create(ctx context.Context, t *models.TeacherSubjects) error
	Delete(ctx context.Context, t *models.TeacherSubjects) error
}

type teacherSubjectRepo struct {
	db *gorm.DB
}

func NewTeacherSubjectRepository(db *gorm.DB) TeacherSubjectRepository {
	return &teacherSubjectRepo{
		db: db,
	}
}

func (r *teacherSubjectRepo) Create(ctx context.Context, t *models.TeacherSubjects) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *teacherSubjectRepo) Delete(ctx context.Context, t *models.TeacherSubjects) error {
	return r.db.WithContext(ctx).Delete(t).Error
}
