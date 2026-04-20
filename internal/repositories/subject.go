// Package repositories
package repositories

import (
	"context"
	"strings"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type SubjectRepository interface {
	Create(ctx context.Context, t *models.Subjects) error
	FilterAndCreateInBatches(ctx context.Context, t []*models.Subjects) error
	Update(ctx context.Context, t *models.Subjects) error
	Delete(ctx context.Context, t *models.Subjects) error
	GetByID(ctx context.Context, id int64) (*models.Subjects, error)
	ExistByModel(ctx context.Context, t *models.Subjects) bool
	List(ctx context.Context, schoolID int64, limit int) ([]*models.Subjects, error)
}

type subjectRepo struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *subjectRepo {
	return &subjectRepo{
		db: db,
	}
}

// Create saves all relational tables
func (r *subjectRepo) Create(ctx context.Context, t *models.Subjects) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *subjectRepo) FilterAndCreateInBatches(
	ctx context.Context, t []*models.Subjects,
) error {
	// 1. Fetch all existing names and codes in two queries (or one combined)
	var existingNames, existingCodes []string
	r.db.WithContext(ctx).Model(&models.Subjects{}).Pluck("name", &existingNames)
	r.db.WithContext(ctx).Model(&models.Subjects{}).Pluck("code", &existingCodes)

	// 2. Convert to maps for O(1) lookup
	nameMap := make(map[string]bool)
	for _, n := range existingNames {
		nameMap[strings.ToLower(n)] = true
	}

	codeMap := make(map[string]bool)
	for _, c := range existingCodes {
		codeMap[strings.ToLower(c)] = true
	}

	// 3. Filter the batch
	var finalBatch []*models.Subjects
	for _, s := range t {
		if !nameMap[strings.ToLower(s.Name)] && !codeMap[strings.ToLower(s.Code)] {
			finalBatch = append(finalBatch, s)
			// Add to maps to prevent duplicates within the SAME batch
			nameMap[strings.ToLower(s.Name)] = true
			codeMap[strings.ToLower(s.Code)] = true
		}
	}

	// 4. Perform a clean batch insert
	if len(finalBatch) > 0 {
		return r.db.WithContext(ctx).CreateInBatches(finalBatch, 100).Error
	}

	return nil
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
