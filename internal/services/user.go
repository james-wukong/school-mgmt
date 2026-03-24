package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type AdminUserService struct {
	repo repositories.AdminUserRepository
}

func NewAdminUserService(db *gorm.DB) *AdminUserService {
	return &AdminUserService{
		repo: repositories.NewAdminUserRepository(db),
	}
}

func (s *AdminUserService) GetUserSchoolID(ctx context.Context, id int64) (*models.AdminUser, error) {
	return s.repo.GetSchoolByUserID(ctx, id)
}
