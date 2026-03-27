// Package repositories
package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(ctx context.Context, t *models.Classes) error
	Update(ctx context.Context, t *models.Classes) error
	UpdateWithSemester(ctx context.Context, t *models.Classes) error
	Delete(ctx context.Context, t *models.Classes) error
	GetByID(ctx context.Context, id int64) (*models.Classes, error)
}

type classRepo struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepo{
		db: db,
	}
}

// Create saves all relational tables
func (r *classRepo) Create(ctx context.Context, t *models.Classes) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// Update will only update classes table
func (r *classRepo) Update(ctx context.Context, t *models.Classes) error {
	return r.db.WithContext(ctx).
		Save(t).
		Error
}

// UpdateWithSemester updates the Class AND upsert the Semesters
func (r *classRepo) UpdateWithSemester(ctx context.Context, t *models.Classes) error {
	// GORM will automatically Save the Semester first,
	// then Save the Class with the correct SemesterID.
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *classRepo) Delete(ctx context.Context, t *models.Classes) error {
	return r.db.WithContext(ctx).Delete(t).Error
}

func (r *classRepo) GetByID(ctx context.Context, id int64) (*models.Classes, error) {
	var s models.Classes

	err := r.db.WithContext(ctx).
		Preload("Semester").
		First(&s, id).
		Error

	if err != nil {
		return nil, err
	}

	return &s, nil
}
