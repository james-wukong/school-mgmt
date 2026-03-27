package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type TeacherRepository interface {
	Create(ctx context.Context, t *models.Teachers) error
	CreateWithTeacherSubject(
		ctx context.Context, t *models.Teachers, ts []*models.TeacherSubjects,
	) error
	Update(ctx context.Context, t *models.Teachers) error
	GetByID(ctx context.Context, id int64) (*models.Teachers, error)
	UpdateTeacherStatus(ctx context.Context, t *models.Teachers) error
	UpdateWithTeacherSubject(
		ctx context.Context, t *models.Teachers, ts []*models.TeacherSubjects,
	) error
	ReplaceWithSubjectAssoc(ctx context.Context, t *models.Teachers) error
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

func (r *teacherRepo) CreateWithTeacherSubject(
	ctx context.Context, t *models.Teachers, ts []*models.TeacherSubjects,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Teacher
		if err := tx.Omit("School", "Subjects").Save(t).Error; err != nil {
			return err
		}

		// 2. Update teacher-subject pairs (teacher id)
		for _, pair := range ts {
			pair.TeacherID = t.ID
		}

		// 3. Insert teacher-subject pairs
		return tx.Model(&models.TeacherSubjects{}).Create(ts).Error
	})
}

func (r *teacherRepo) Update(ctx context.Context, t *models.Teachers) error {
	// r.db.Model(&t).Select("IsActive").Updates(t).Error
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *teacherRepo) GetByID(ctx context.Context, id int64) (*models.Teachers, error) {
	var teacher models.Teachers

	err := r.db.WithContext(ctx).
		Preload("School").
		Preload("Subjects").
		First(&teacher, id).Error

	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

// UpdateTeacherStatus only update is_active field for a teacher
func (r *teacherRepo) UpdateTeacherStatus(
	ctx context.Context, t *models.Teachers,
) error {
	// r.db.Model(&t).Select("IsActive").Updates(t).Error
	return r.db.WithContext(ctx).
		Model(t).
		Where("id = ?", t.ID).
		UpdateColumn("is_active", t.IsActive).
		Error
}

// UpdateWithTeacherSubject update teachers table
// removes previously attached teacher-subject pairs
// and insert new pairs
func (r *teacherRepo) UpdateWithTeacherSubject(
	ctx context.Context, t *models.Teachers, ts []*models.TeacherSubjects,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Teacher
		if err := tx.Omit("School", "Subjects").Save(t).Error; err != nil {
			return err
		}

		// 2. Remove the previous teacher-subject pairs
		if err := tx.Where("teacher_id = ?", t.ID).
			Delete(&models.TeacherSubjects{}).Error; err != nil {
			return err
		}
		// 3. Insert teacher-subject pairs
		return tx.Model(&models.TeacherSubjects{}).Create(ts).Error
	})
}

// ReplaceWithSubjectAssoc save a teacher model,
// upsert subjects associated with it,
// and remove all other previously associated subjects
func (r *teacherRepo) ReplaceWithSubjectAssoc(ctx context.Context, t *models.Teachers) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Teacher
		if err := tx.Omit("School", "Subjects").Save(t).Error; err != nil {
			return err
		}

		// 2. Replace the association
		// This Upserts the current slice and Deletes/Unlinks others
		return tx.Model(t).
			Unscoped().
			Association("Subjects").
			Replace(t.Subjects)
	})
}
