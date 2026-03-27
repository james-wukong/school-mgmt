package services

import (
	"context"
	"errors"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type TeacherService struct {
	repo repositories.TeacherRepository
}

func NewTeacherService(db *gorm.DB) *TeacherService {
	return &TeacherService{
		repo: repositories.NewTeacherRepository(db),
	}
}

func (s *TeacherService) GetTeacher(ctx context.Context, id int64) (*models.Teachers, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TeacherService) CreateTeacher(ctx context.Context, t *models.Teachers) error {
	if t.FirstName == "" {
		return errors.New("first name required")
	}

	return s.repo.Create(ctx, t)
}

func (s *TeacherService) CreateWithTeacherSubject(
	ctx context.Context, t *models.Teachers, ts []*models.TeacherSubjects,
) error {

	return s.repo.CreateWithTeacherSubject(ctx, t, ts)
}

func (s *TeacherService) UpdateStatus(ctx context.Context, t *models.Teachers) error {
	return s.repo.UpdateTeacherStatus(ctx, t)
}

func (s *TeacherService) UpdateWithAssoc(ctx context.Context, t *models.Teachers) error {
	return s.repo.Update(ctx, t)
}

func (s *TeacherService) ReplaceWithSubjectAssoc(ctx context.Context, t *models.Teachers) error {
	return s.repo.ReplaceWithSubjectAssoc(ctx, t)
}

func (s *TeacherService) UpdateWithTeacherSubject(
	ctx context.Context, t *models.Teachers, ts []*models.TeacherSubjects,
) error {
	return s.repo.UpdateWithTeacherSubject(ctx, t, ts)
}
