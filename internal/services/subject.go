package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type SubjectService struct {
	repo repositories.SubjectRepository
}

func NewSubjectService(db *gorm.DB) *SubjectService {
	return &SubjectService{
		repo: repositories.NewSubjectRepository(db),
	}
}

func (s *SubjectService) CreateSubject(ctx context.Context, t *models.Subjects) error {
	return s.repo.Create(ctx, t)
}

func (s *SubjectService) GetByID(ctx context.Context, subjectID int64) (*models.Subjects, error) {
	return s.repo.GetByID(ctx, subjectID)
}

func (s *SubjectService) List(
	ctx context.Context,
	schoolID int64,
	limit int,
) ([]*models.Subjects, error) {
	return s.repo.List(ctx, schoolID, limit)
}
