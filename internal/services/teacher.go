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

func (s *TeacherService) CreateWithAssoc(
	ctx context.Context, t *models.Teachers,
	ts []*models.TeacherSubjects, tt []*models.TeacherTimeslots,
) error {
	return s.repo.CreateWithAssoc(ctx, t, ts, tt)
}

func (s *TeacherService) CreateWithTeacherTimeslot(
	ctx context.Context, t *models.Teachers, tt []*models.TeacherTimeslots,
) error {
	return s.repo.CreateWithTeacherTimeslot(ctx, t, tt)
}

func (s *TeacherService) CreateWithSubjectJoinInBatches(
	ctx context.Context, t []*models.Teachers,
) error {
	return s.repo.CreateWithSubjectJoinInBatches(ctx, t)
}

func (s *TeacherService) UpdateStatus(ctx context.Context, t *models.Teachers) error {
	return s.repo.UpdateTeacherStatus(ctx, t)
}

func (s *TeacherService) UpdateWithAssoc(
	ctx context.Context,
	t *models.Teachers,
	ts []*models.TeacherSubjects,
	tt []*models.TeacherTimeslots,
	semID int64,
) error {
	return s.repo.UpdateWithAssoc(ctx, t, ts, tt, semID)
}

func (s *TeacherService) ReplaceWithSubjectAssoc(ctx context.Context, t *models.Teachers) error {
	return s.repo.ReplaceWithSubjectAssoc(ctx, t)
}
