package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TimeslotRepository interface {
	Create(ctx context.Context, t *models.Timeslots) error
	CreateInBatches(ctx context.Context, t []*models.Timeslots) error
	Update(ctx context.Context, t *models.Timeslots) error
	GetByID(ctx context.Context, id int64) (*models.Timeslots, error)
	UpdateWithAssoc(ctx context.Context, t *models.Timeslots) error
	Delete(ctx context.Context, t *models.Timeslots) error
	List(ctx context.Context, schoolID int64, semester int64, limit int) ([]*models.Timeslots, error)
}

type timeslotsRepo struct {
	db *gorm.DB
}

func NewTimeslotRepository(db *gorm.DB) *timeslotsRepo {
	return &timeslotsRepo{
		db: db,
	}
}

func (r *timeslotsRepo) Create(ctx context.Context, t *models.Timeslots) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *timeslotsRepo) CreateInBatches(ctx context.Context, t []*models.Timeslots) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Omit("Semester", "School").
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "semester_id"},
					{Name: "day_of_week"},
					{Name: "start_time"},
				},
				DoNothing: true,
			}).
			CreateInBatches(t, 100).
			Error
	})
}

func (r *timeslotsRepo) Update(ctx context.Context, t *models.Timeslots) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *timeslotsRepo) GetByID(ctx context.Context, id int64) (*models.Timeslots, error) {
	var u models.Timeslots

	err := r.db.
		WithContext(ctx).
		First(&u, id).Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *timeslotsRepo) UpdateWithAssoc(ctx context.Context, t *models.Timeslots) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *timeslotsRepo) Delete(ctx context.Context, t *models.Timeslots) error {
	return r.db.WithContext(ctx).Delete(t).Error
}

// List timeslots by school id and semester id
func (r *timeslotsRepo) List(
	ctx context.Context, schoolID int64, semesterID int64, limit int,
) ([]*models.Timeslots, error) {
	var l []*models.Timeslots
	sql := r.db.WithContext(ctx).
		Where("school_id = ? AND semester_id = ?", schoolID, semesterID)
	if limit != 0 {
		sql = sql.Limit(limit)
	}
	err := sql.
		Order("day_of_week ASC").
		Find(&l).
		Error

	if err != nil {
		return nil, err
	}
	return l, nil
}
