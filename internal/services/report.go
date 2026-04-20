package services

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type ReportService struct {
	repo repositories.ReportRepository
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{
		repo: repositories.NewReportRepository(db),
	}
}
func (s *ReportService) GetWeeklyClassReport(
	ctx context.Context, semesterID int64, version float64,
) ([]models.WeeklyClassScheduleReport, error) {
	return s.repo.GetWeeklyClassReport(ctx, semesterID, version)
}

func (s *ReportService) GetWeeklyTeacherReport(
	ctx context.Context, semesterID int64, version float64,
) ([]models.WeeklyTeacherScheduleReport, error) {
	return s.repo.GetWeeklyTeacherReport(ctx, semesterID, version)
}

func (s *ReportService) GetMaxDay(
	ctx context.Context, semesterID int64, version float64,
) int {
	return s.repo.GetMaxDay(ctx, semesterID, version)
}
