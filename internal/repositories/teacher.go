package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type TeacherRepository interface {
	Create(ctx context.Context, t *models.Teachers) error
	CreateWithAssoc(
		ctx context.Context, t *models.Teachers,
		ts []*models.TeacherSubjects,
		tt []*models.TeacherTimeslots,
	) error
	CreateWithTeacherTimeslot(
		ctx context.Context, t *models.Teachers, tt []*models.TeacherTimeslots,
	) error
	Update(ctx context.Context, t *models.Teachers) error
	GetByID(ctx context.Context, id int64) (*models.Teachers, error)
	UpdateTeacherStatus(ctx context.Context, t *models.Teachers) error
	UpdateWithAssoc(
		ctx context.Context, t *models.Teachers,
		ts []*models.TeacherSubjects,
		tt []*models.TeacherTimeslots,
		semID int64,
	) error
	ReplaceWithSubjectAssoc(ctx context.Context, t *models.Teachers) error
	ExistByModel(ctx context.Context, t *models.Teachers) bool
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

func (r *teacherRepo) CreateWithAssoc(
	ctx context.Context, t *models.Teachers,
	ts []*models.TeacherSubjects, tt []*models.TeacherTimeslots,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create the Teacher
		if err := tx.Omit("School", "Subjects", "Timeslots").Create(t).Error; err != nil {
			return err
		}

		// 2. Update teacher-subject pairs (teacher id)
		for _, spair := range ts {
			spair.TeacherID = t.ID
		}
		for _, tpair := range tt {
			tpair.TeacherID = t.ID
		}

		// 3. Insert teacher-subject pairs
		if len(ts) > 0 {
			if err := tx.Model(&models.TeacherSubjects{}).Create(ts).Error; err != nil {
				return err
			}
		}
		if len(tt) > 0 {
			if err := tx.Model(&models.TeacherTimeslots{}).Create(tt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *teacherRepo) CreateWithTeacherTimeslot(
	ctx context.Context, t *models.Teachers, tt []*models.TeacherTimeslots,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Teacher
		if err := tx.Omit("School", "Subjects", "Timeslots").Save(t).Error; err != nil {
			return err
		}

		// 2. Update teacher-timeslot pairs (teacher id)
		for _, pair := range tt {
			pair.TeacherID = t.ID
		}

		// 3. Insert teacher-timeslot pairs
		if len(tt) > 0 {
			return tx.Model(&models.TeacherTimeslots{}).Create(tt).Error
		}
		return nil
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
		Preload("Timeslots").
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

// UpdateWithAssoc update teachers table
// removes previously attached teacher-subject and teacher-timeslot pairs
// and insert new pairs
func (r *teacherRepo) UpdateWithAssoc(
	ctx context.Context, t *models.Teachers,
	ts []*models.TeacherSubjects,
	tt []*models.TeacherTimeslots,
	semID int64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Teacher
		if err := tx.Omit("School", "Subjects", "Timeslots").Save(t).Error; err != nil {
			return err
		}

		// 2.1 Remove the previous teacher-subject pairs
		if err := tx.Where("teacher_id = ?", t.ID).
			Delete(&models.TeacherSubjects{}).Error; err != nil {
			return err
		}
		// 2.2 Remove the previous teacher-timeslot pair for the semester
		if semID != 0 {
			subQuery := tx.Model(&models.Timeslots{}).
				Select("id").
				Where("semester_id = ?", semID)
			if err := tx.Where("teacher_id = ?", t.ID).
				Where("timeslot_id IN (?)", subQuery).
				Delete(&models.TeacherTimeslots{}).
				Error; err != nil {
				return err
			}
		}

		// 3. Insert teacher-subject pairs
		if len(ts) > 0 {
			if err := tx.Model(&models.TeacherSubjects{}).Create(ts).Error; err != nil {
				return err
			}
		}

		if len(tt) > 0 {
			if err := tx.Model(&models.TeacherTimeslots{}).Create(tt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ReplaceWithSubjectAssoc save a teacher model,
// upsert subjects associated with it,
// and remove all other previously associated subjects
func (r *teacherRepo) ReplaceWithSubjectAssoc(ctx context.Context, t *models.Teachers) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Teacher
		if err := tx.Omit("School", "Subjects", "Timeslots").Save(t).Error; err != nil {
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

func (r *teacherRepo) ExistByModel(ctx context.Context, t *models.Teachers) bool {
	var exists bool
	r.db.WithContext(ctx).
		Model(&models.Teachers{}).
		Select("count(*) > 0").
		Where(t).
		Limit(1).
		Find(&exists)

	return exists
}
