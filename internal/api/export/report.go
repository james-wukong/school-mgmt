package export

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/james-wukong/online-school-mgmt/internal/config"
	"github.com/james-wukong/online-school-mgmt/internal/dto"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	utils "github.com/james-wukong/online-school-mgmt/internal/utils/export"
	"gorm.io/gorm"
)

func ExportReportHandler(dbConn *gorm.DB) context.Handler {
	return func(ctx *context.Context) {
		var reqData dto.ScheduleExportRequest

		// 1. Decode the body directly into the struct
		err := json.NewDecoder(ctx.Request.Body).Decode(&reqData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"msg":     err.Error(),
				"success": false,
			})
			return
		}
		// 2. Load extra configurations
		cfg := config.InitConfig()
		// Save to csv file
		classFileName := fmt.Sprintf("class_report_s_%d_v_%.2f.cvs",
			reqData.SemesterID, reqData.SchedVersion.InexactFloat64(),
		)
		teacherFileName := fmt.Sprintf("teacher_report_s_%d_v_%.2f.cvs",
			reqData.SemesterID, reqData.SchedVersion.InexactFloat64(),
		)

		reportRepo := repositories.NewReportRepository(dbConn)
		classService := utils.NewClassReportService(reportRepo)
		teacherService := utils.NewTeacherReportService(reportRepo)
		var reportService repositories.ReportService
		for _, filename := range []string{classFileName, teacherFileName} {
			f, err := os.Create(filepath.Join(
				cfg.App.ExportDownloadURI, filepath.Base(filename),
			))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, map[string]interface{}{
					"code":    http.StatusBadRequest,
					"msg":     err.Error(),
					"success": false,
				})
				return
			}
			defer f.Close()
			// Write the UTF-8 BOM bytes first
			if _, err := f.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
				ctx.JSON(http.StatusBadRequest, map[string]interface{}{
					"code":    http.StatusBadRequest,
					"msg":     err.Error(),
					"success": false,
				})
				return
			}
			switch filename {
			case classFileName:
				reportService = classService
			case teacherFileName:
				reportService = teacherService
			}
			if err := reportService.ExportToCSV(ctx.Request.Context(),
				f, reqData.SemesterID, reqData.SchedVersion.InexactFloat64(),
			); err != nil {
				ctx.JSON(http.StatusBadRequest, map[string]interface{}{
					"code":    http.StatusBadRequest,
					"msg":     err.Error(),
					"success": false,
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"msg":     "Success",
			"success": true,
		})
	}
}
