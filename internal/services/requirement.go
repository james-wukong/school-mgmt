package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RequirementService struct {
	repo repositories.RequirementRepository
}

func NewRequirementService(db *gorm.DB) *RequirementService {
	return &RequirementService{
		repo: repositories.NewRequirementRepository(db),
	}
}

func (s *RequirementService) CreateWithAssocInBatch(ctx context.Context, t []*models.Requirements) error {
	return s.repo.CreateWithAssocInBatch(ctx, t)
}

func (s *RequirementService) GetNewVersion(ctx context.Context, semesterID int64) decimal.Decimal {
	return s.repo.GetVersion(ctx, semesterID)
}
