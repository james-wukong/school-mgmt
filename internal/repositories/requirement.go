package repositories

import (
	"context"
	"errors"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequirementRepository interface {
	CreateWithAssoc(ctx context.Context, t *models.Requirements) error
	SaveInBatch(ctx context.Context, t []*models.Requirements) error
	GetByID(ctx context.Context, id int64) (*models.Requirements, error)
	UpdateWithAssoc(ctx context.Context, t *models.Requirements) error
	Delete(ctx context.Context, t *models.Requirements) error
	ListBySemVersion(
		ctx context.Context, semesterID int64, version float64,
	) ([]*models.Requirements, error)
	GetVersion(ctx context.Context, semesterID int64) decimal.Decimal
}

type requirementsRepo struct {
	db *gorm.DB
}

func NewRequirementRepository(db *gorm.DB) RequirementRepository {
	return &requirementsRepo{
		db: db,
	}
}

func (r *requirementsRepo) CreateWithAssoc(ctx context.Context, t *models.Requirements) error {
	// 1. Inserts Requirements if not exists
	// 2. Updates Requirements if conflict occurs
	// 3. Inserts/updates associated records
	// 4. Updates join table (many2many) entries if there is any
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Session(&gorm.Session{
			FullSaveAssociations: true,
		}).Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "school_id"},
				{Name: "semester_id"},
				{Name: "subject_id"},
				{Name: "teacher_id"},
				{Name: "class_id"},
				{Name: "version"},
			}, // or unique key
			UpdateAll: true,
		}).Create(t).
			Error
	})
}

func (r *requirementsRepo) SaveInBatch(ctx context.Context, t []*models.Requirements) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.
			Omit(clause.Associations).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "semester_id"},
					{Name: "subject_id"},
					{Name: "teacher_id"},
					{Name: "class_id"},
					{Name: "version"},
				}, // or unique key
				UpdateAll: true,
			}).CreateInBatches(t, 100).
			Error
	})
}

func (r *requirementsRepo) GetByID(ctx context.Context, id int64) (*models.Requirements, error) {
	var u models.Requirements

	err := r.db.
		WithContext(ctx).
		Preload("Class").
		Preload("School").
		Preload("Semester").
		Preload("Subject").
		Preload("Teacher").
		First(&u, id).Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *requirementsRepo) UpdateWithAssoc(ctx context.Context, t *models.Requirements) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *requirementsRepo) Delete(ctx context.Context, t *models.Requirements) error {
	return r.db.WithContext(ctx).Delete(t).Error
}

// List requirements by school id, semester id and version
func (r *requirementsRepo) ListBySemVersion(
	ctx context.Context,
	semesterID int64, version float64,
) ([]*models.Requirements, error) {
	var l []*models.Requirements
	err := r.db.WithContext(ctx).
		Where("semester_id = ? and version = ?",
			semesterID, version).
		Order("id DESC").
		Find(&l).
		Error

	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *requirementsRepo) GetVersion(ctx context.Context, semesterID int64) decimal.Decimal {
	var m models.Requirements
	err := r.db.WithContext(ctx).
		Select("id", "version").
		Where("semester_id = ?", semesterID).
		Order("version DESC").
		First(&m, 1).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return decimal.NewFromFloat(1.00)
		}
	}

	return m.Version.Add(decimal.NewFromFloat(0.01))
}
