package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type ClassService struct {
	repo repositories.ClassRepository
}

func NewClassService(db *gorm.DB) *ClassService {
	return &ClassService{
		repo: repositories.NewClassRepository(db),
	}
}

func (s *ClassService) GetByID(ctx context.Context, classID int64) (*models.Classes, error) {
	return s.repo.GetByID(ctx, classID)
}
