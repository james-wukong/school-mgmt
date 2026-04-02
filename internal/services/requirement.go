package services

import (
	"context"
	"fmt"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RequirementService struct {
	repo    repositories.RequirementRepository
	clsRepo repositories.ClassRepository
	subRepo repositories.SubjectRepository
	tchRepo repositories.TeacherRepository
}

func NewRequirementService(db *gorm.DB) *RequirementService {
	return &RequirementService{
		repo:    repositories.NewRequirementRepository(db),
		clsRepo: repositories.NewClassRepository(db),
		subRepo: repositories.NewSubjectRepository(db),
		tchRepo: repositories.NewTeacherRepository(db),
	}
}

func (s *RequirementService) SaveRequirements(ctx context.Context, t []*models.Requirements) error {
	return s.repo.SaveInBatch(ctx, t)
}

func (s *RequirementService) GetNewVersion(ctx context.Context, semesterID int64) decimal.Decimal {
	return s.repo.GetVersion(ctx, semesterID)
}

func (s *RequirementService) ValidateAssoc(
	ctx context.Context, t *models.Requirements,
) []error {
	var errs []error
	if exists := s.subRepo.ExistByModel(ctx, t.Subject); !exists {
		errs = append(errs, fmt.Errorf("subject: id %d name %s does not exist",
			t.Subject.ID, t.Subject.Name,
		))
	}
	if exists := s.tchRepo.ExistByModel(ctx, t.Teacher); !exists {
		errs = append(errs, fmt.Errorf("teacher: id %d first name %s last name: %s does not exist",
			t.Teacher.ID, t.Teacher.FirstName, t.Teacher.LastName,
		))
	}
	if exists := s.clsRepo.ExistByModel(ctx, t.Class); !exists {
		errs = append(errs, fmt.Errorf("class: id %d grade %d class: %s does not exist",
			t.Class.ID, t.Class.Grade, t.Class.ClassName,
		))
	}
	return errs
}
