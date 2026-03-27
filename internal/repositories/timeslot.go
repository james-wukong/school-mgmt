package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type TimeslotRepository interface {
	Create(ctx context.Context, t *models.Timeslots) error
	Update(ctx context.Context, t *models.Timeslots) error
	GetByID(ctx context.Context, id int64) (*models.Timeslots, error)
	UpdateWithAssoc(ctx context.Context, t *models.Timeslots) error
	Delete(ctx context.Context, t *models.Timeslots) error
}

type timeslotsRepo struct {
	db *gorm.DB
}

func NewTimeslotRepository(db *gorm.DB) TimeslotRepository {
	return &timeslotsRepo{
		db: db,
	}
}

func (r *timeslotsRepo) Create(ctx context.Context, t *models.Timeslots) error {
	return r.db.WithContext(ctx).Create(t).Error
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
