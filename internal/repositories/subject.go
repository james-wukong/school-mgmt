// Package repositories
package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type SubjectRepository interface {
	Create(ctx context.Context, t *models.Subjects) error
	Update(ctx context.Context, t *models.Subjects) error
	Delete(ctx context.Context, t *models.Subjects) error
	GetByID(ctx context.Context, id int64) (*models.Subjects, error)
	ExistByModel(ctx context.Context, t *models.Subjects) bool
	List(ctx context.Context, schoolID int64, limit int) ([]*models.Subjects, error)
}

type subjectRepo struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepo{
		db: db,
	}
}

// Create saves all relational tables
func (r *subjectRepo) Create(ctx context.Context, t *models.Subjects) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// Update will only update classes table
func (r *subjectRepo) Update(ctx context.Context, t *models.Subjects) error {
	return r.db.WithContext(ctx).
		Save(t).
		Error
}

func (r *subjectRepo) Delete(ctx context.Context, t *models.Subjects) error {
	return r.db.WithContext(ctx).Delete(t).Error
}

func (r *subjectRepo) GetByID(ctx context.Context, id int64) (*models.Subjects, error) {
	var s models.Subjects

	err := r.db.WithContext(ctx).
		Preload("School").
		Preload("Teachers").
		First(&s, id).
		Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *subjectRepo) ExistByModel(ctx context.Context, t *models.Subjects) bool {
	var exists bool
	r.db.WithContext(ctx).
		Model(&models.Subjects{}).
		Select("count(*) > 0").
		Where(t).
		Limit(1).
		Find(&exists)

	return exists
}

func (r *subjectRepo) List(
	ctx context.Context, schoolID int64, limit int,
) ([]*models.Subjects, error) {
	var l []*models.Subjects
	sql := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID)
	if limit != 0 {
		sql = sql.Limit(limit)
	}
	err := sql.
		Order("id DESC").
		Find(&l).
		Error

	if err != nil {
		return nil, err
	}
	return l, nil
}
