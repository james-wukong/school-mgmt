package repositories

import (
	"context"
	"io"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportRepository interface {
	GetWeeklyClassReport(
		ctx context.Context, semesterID int64, version decimal.Decimal,
	) ([]models.WeeklyClassScheduleReport, error)

	GetWeeklyTeacherReport(
		ctx context.Context, semesterID int64, version decimal.Decimal,
	) ([]models.WeeklyTeacherScheduleReport, error)

	GetMaxDay(
		ctx context.Context, semesterID int64, version decimal.Decimal,
	) int
}

type ReportService interface {
	ExportToCSV(ctx context.Context, w io.Writer, semesterID int64, version decimal.Decimal) error
}

type reportsRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *reportsRepo {
	return &reportsRepo{
		db: db,
	}
}

func (r *reportsRepo) GetWeeklyClassReport(
	ctx context.Context, semesterID int64, version decimal.Decimal,
) ([]models.WeeklyClassScheduleReport, error) {
	// get report from VIEW: vw_class_weekly_schedule
	var schedules []models.WeeklyClassScheduleReport

	// GORM will execute:
	// SELECT * FROM v_weekly_schedules ORDER BY grade, class_name, start_time ASC
	if err := r.db.WithContext(ctx).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "grade"}, Desc: false},
			{Column: clause.Column{Name: "class_name"}, Desc: false},
			{Column: clause.Column{Name: "start_time"}, Desc: false},
			{Column: clause.Column{Name: "day_of_week"}, Desc: false},
		}}).
		Find(&schedules, "semester_id = ? AND version = ?", semesterID, version).
		Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *reportsRepo) GetWeeklyTeacherReport(
	ctx context.Context, semesterID int64, version decimal.Decimal,
) ([]models.WeeklyTeacherScheduleReport, error) {
	// get report from VIEW: vw_teacher_weekly_schedule
	var schedules []models.WeeklyTeacherScheduleReport

	// GORM will execute:
	// SELECT * FROM v_weekly_schedules ORDER BY grade, class_name, start_time ASC
	if err := r.db.WithContext(ctx).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "teacher_id"}, Desc: false},
			{Column: clause.Column{Name: "start_time"}, Desc: false},
			{Column: clause.Column{Name: "day_of_week"}, Desc: false},
		}}).
		Find(&schedules, "semester_id = ? AND version = ?", semesterID, version).
		Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *reportsRepo) GetMaxDay(
	ctx context.Context, semesterID int64, version decimal.Decimal,
) int {
	var maxDay int
	if err := r.db.WithContext(ctx).
		Model(&models.WeeklyClassScheduleReport{}).
		Select("MAX(day_of_week)").
		Where("semester_id = ? AND version = ?", semesterID, version).
		Row().
		Scan(&maxDay); err != nil {
		return 0
	}
	return maxDay
}
