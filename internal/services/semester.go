// Package services
package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type SemesterService struct {
	repo repositories.SemesterRepository
}

func NewSemesterService(db *gorm.DB) *SemesterService {
	return &SemesterService{
		repo: repositories.NewSemesterRepository(db),
	}
}

func (s *SemesterService) List(
	ctx context.Context,
	schoolID int64,
	limit int,
) ([]*models.Semesters, error) {
	return s.repo.List(ctx, schoolID, limit)
}

func (s *SemesterService) SaveWithAssoc(ctx context.Context, t *models.Semesters) error {
	if len(t.Classes) == 0 {
		return models.ErrEmptyAssociations
	}
	return s.repo.UpdateWithAssoc(ctx, t)
}

func (s *SemesterService) GetByID(ctx context.Context, id int64) (*models.Semesters, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SemesterService) AppendClasses(
	ctx context.Context, t *models.Semesters, c []*models.Classes,
) error {
	return s.repo.AppendClasses(ctx, t, c)
}

func (s *SemesterService) ReplaceWithTimeslotAssoc(
	ctx context.Context, t *models.Semesters,
) error {
	return s.repo.ReplaceWithTimeslotAssoc(ctx, t)
}
