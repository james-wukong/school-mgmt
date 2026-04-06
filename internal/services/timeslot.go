package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type TimeslotService struct {
	repo repositories.TimeslotRepository
}

func NewTimeslotService(db *gorm.DB) *TimeslotService {
	return &TimeslotService{
		repo: repositories.NewTimeslotRepository(db),
	}
}

func (s *TimeslotService) CreateTimeslotsInBatches(ctx context.Context, t []*models.Timeslots) error {
	// return s.repo.Create(ctx, t)
	return s.repo.CreateInBatches(ctx, t)
}

func (s *TimeslotService) List(
	ctx context.Context,
	schoolID int64,
	semesterID int64,
	limit int,
) ([]*models.Timeslots, error) {
	return s.repo.List(ctx, schoolID, semesterID, limit)
}
