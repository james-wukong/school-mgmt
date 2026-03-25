package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type SemesterRepository interface {
	Create(ctx context.Context, t *models.Semesters) error
	UpdateWithClasses(ctx context.Context, t *models.Semesters) error
	UpdateWithClassAssocReplace(ctx context.Context, t *models.Semesters) error
	AppendClasses(ctx context.Context, t *models.Semesters, c []*models.Classes) error
	Delete(ctx context.Context, t *models.Semesters) error
	GetByID(ctx context.Context, id int64) (*models.Semesters, error)
	List(ctx context.Context, schoolID int64, limit int) ([]*models.Semesters, error)
}

type semesterRepo struct {
	db *gorm.DB
}

func NewSemesterRepository(db *gorm.DB) SemesterRepository {
	return &semesterRepo{
		db: db,
	}
}

// Create will insert all association classes into database
func (r *semesterRepo) Create(ctx context.Context, t *models.Semesters) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *semesterRepo) UpdateWithClasses(ctx context.Context, t *models.Semesters) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *semesterRepo) UpdateWithClassAssocReplace(ctx context.Context, t *models.Semesters) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Parent
		if err := tx.Omit("Classes").Save(t).Error; err != nil {
			return err
		}

		// 2. Replace the association
		// This Upserts the current slice and Deletes/Unlinks others
		return tx.Model(t).Association("Classes").Replace(t.Classes)
	})
}

func (r *semesterRepo) AppendClasses(
	ctx context.Context,
	t *models.Semesters,
	c []*models.Classes,
) error {
	// This only executes the INSERT for the new classes
	return r.db.
		WithContext(ctx).
		Model(t).
		Association("Classes").
		Append(c)
}

func (r *semesterRepo) Delete(ctx context.Context, t *models.Semesters) error {
	return r.db.WithContext(ctx).Delete(t).Error
}

func (r *semesterRepo) GetByID(ctx context.Context, id int64) (*models.Semesters, error) {
	var s models.Semesters

	err := r.db.WithContext(ctx).
		Preload("Classes").
		Preload("School").
		First(&s, id).
		Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *semesterRepo) List(
	ctx context.Context, schoolID int64, limit int,
) ([]*models.Semesters, error) {
	var l []*models.Semesters
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Limit(limit).
		Order("id DESC").
		Find(&l).
		Error

	if err != nil {
		return nil, err
	}
	return l, nil
}
