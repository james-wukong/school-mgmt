// Package export mainly exports weekly reports in csv files
package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"github.com/james-wukong/online-school-mgmt/internal/types"
	"github.com/shopspring/decimal"
)

// ClassReportService handles the export logic
type ClassReportService struct {
	repo repositories.ReportRepository
}

// NewClassReportService creates a new instance
func NewClassReportService(repo repositories.ReportRepository) *ClassReportService {
	return &ClassReportService{repo: repo}
}

// ExportToCSV writes the schedule data to the provided writer
func (s *ClassReportService) ExportToCSV(
	ctx context.Context, w io.Writer, semesterID int64, version decimal.Decimal,
) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	// 1. Fetch sorted data
	rows, err := s.repo.GetWeeklyClassReport(ctx, semesterID, version)
	if err != nil {
		return err
	}
	// 2. Get max week day
	maxDay := s.repo.GetMaxDay(ctx, semesterID, version)

	// 3. Write Headers
	headers := createHeaders(maxDay)
	if err := writeHeaders(writer, headers); err != nil {
		return err
	}

	// 4. Process Rows
	var currentClass string
	var currentStartTime types.ClockTime
	displayRow := make([]string, len(headers))
	for i := range rows {
		if !time.Time(currentStartTime).IsZero() && currentStartTime != rows[i].StartTime {
			// Insert row
			if err := writer.Write(displayRow); err != nil {
				return err
			}
			displayRow = make([]string, len(headers))
		}
		// Insert break line on class change
		if currentClass != "" && currentClass != rows[i].ClassName {
			if err := writer.Write([]string{}); err != nil {
				return err
			}
		}
		// fill display slice
		if displayRow[0] == "" {
			displayRow[0] = fmt.Sprintf("%d - (%s)", rows[i].Grade, rows[i].ClassName)
		}
		if displayRow[1] == "" {
			displayRow[1] = time.Time(rows[i].StartTime).Format("15:04")
		}
		displayRow[int(rows[i].DayOfWeek)+1] = fmt.Sprintf("%s-%s",
			rows[i].TeacherName,
			rows[i].SubjectName,
		)

		currentClass = rows[i].ClassName
		currentStartTime = rows[i].StartTime
	}

	// After the loop ends, check if there is an unwritten row
	// This writes the last displayRow to CSV file
	if !time.Time(currentStartTime).IsZero() {
		if err := writer.Write(displayRow); err != nil {
			return err
		}
	}
	writer.Flush() // Always a good habit to flush the buffer at the end
	return nil
}

func createHeaders(maxDay int) []string {
	dayIndex := map[int]string{
		2: "Monday", 3: "Tuesday", 4: "Wednesday",
		5: "Thursday", 6: "Friday", 7: "Saturday", 8: "Sunday",
	}
	headers := []string{"Class", "Timeslot"}
	for day := range maxDay {
		if day <= 6 {
			headers = append(headers, dayIndex[day+2])
		}
	}
	return headers
}

func writeHeaders(w *csv.Writer, headers []string) error {
	if err := w.Write(headers); err != nil {
		return err
	}
	return nil
}
