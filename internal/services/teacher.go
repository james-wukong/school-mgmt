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
	teacher, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !teacher.IsActive {
		return nil, errors.New("teacher is inactive")
	}

	return teacher, nil
}

func (s *TeacherService) CreateTeacher(ctx context.Context, t *models.Teachers) error {
	if t.FirstName == "" {
		return errors.New("first name required")
	}

	return s.repo.Create(ctx, t)
}
