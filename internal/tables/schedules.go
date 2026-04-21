package tables

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/config"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetSchedulesTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		schedules := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		cfg := config.InitConfig()
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		semService := services.NewSemesterService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		semesters, err := semService.List(ctx.Request.Context(), u.SchoolID, 6)
		if err != nil {
			panic(err)
		}
		var semFilterOpts types.FieldOptions
		for _, semester := range semesters {
			semFilterOpts = append(semFilterOpts, types.FieldOption{
				Value: fmt.Sprint(semester.ID),
				Text: fmt.Sprintf("ID: %d - Year: %d - Semester: %d",
					semester.ID, semester.Year, semester.Semester,
				),
			})
		}
		info := schedules.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}
		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Semester", "semester_id", db.Int8).
			FieldDisplay(func(model types.FieldModel) interface{} {
				semID, err := strconv.ParseInt(model.Value, 10, 64)
				if err != nil {
					panic(err)
				}
				sem, err := semService.GetByID(ctx.Request.Context(), semID)
				if err != nil {
					panic(err)
				}
				return fmt.Sprintf("ID: %s - Year: %d - Semester: %d",
					model.Value, sem.Year, sem.Semester,
				)
			}).
			FieldFilterable(types.FilterType{
				FormType: form.SelectSingle,
				Width:    450,
			}).
			FieldFilterOptions(semFilterOpts).
			FieldSortable()

		info.AddField("Schedule Version", "version", db.Varchar)
		info.AddField("Requirement Version", "requirement_version", db.Varchar)
		info.AddField("Status", "status", db.Varchar)
		info.AddField("Class Report", "class_report", db.Varchar).
			FieldDisplay(func(model types.FieldModel) interface{} {
				fname := printExportFilename("class_report",
					model.Row["version"].(string),
					model.Row["semester_id"].(int64),
				)
				fname = filepath.Join(cfg.App.ExportDownloadPath, fname)
				// check file existence
				if _, err := os.Stat(fname); err != nil {
					fmt.Printf("filename: %s, error: %+v\n", fname, err)
					return ""
				}
				downloadURL, err := url.JoinPath(cfg.App.ExportDownloadURI, filepath.Base(fname))
				if err != nil {
					fmt.Printf("downloadUrl: %s, error: %+v\n", downloadURL, err)
					return ""
				}

				// Return a styled link or button
				return template.HTML(fmt.Sprintf(`
                <a href="%s" target="_blank" class="btn btn-xs btn-primary" download>
                    <i class="fa fa-download"></i> %s
                </a>`, downloadURL, filepath.Base(fname)))
			})

		info.AddField("Teacher Report", "teacher_report", db.Varchar).
			FieldDisplay(func(model types.FieldModel) interface{} {
				fname := printExportFilename("teacher_report",
					model.Row["version"].(string),
					model.Row["semester_id"].(int64),
				)
				fname = filepath.Join(cfg.App.ExportDownloadPath, fname)
				// check file existence
				if _, err := os.Stat(fname); err != nil {
					fmt.Printf("filename: %s, error: %+v\n", fname, err)
					return ""
				}
				downloadURL, err := url.JoinPath(cfg.App.ExportDownloadURI, filepath.Base(fname))
				if err != nil {
					fmt.Printf("downloadUrl: %s, error: %+v\n", downloadURL, err)
					return ""
				}

				// Return a styled link or button
				return template.HTML(fmt.Sprintf(`
                <a href="%s" target="_blank" class="btn btn-xs btn-primary" download>
                    <i class="fa fa-download"></i> %s
                </a>`, downloadURL, filepath.Base(fname)))
			})

		info.HideDeleteButton()
		info.SetTable("schedules").SetTitle("Schedules").SetDescription("Schedules")
		info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			var reqs []*models.ScheduleVersion
			// 1. Create the initial query
			query := dbConn.WithContext(ctx.Request.Context()).
				Table("schedules").
				Select(`DISTINCT ON (semester_id, version)
					schedules.id as id,
					schedules.status as status,
					schedules.semester_id as semester_id,
					schedules.version as version,
					requirements.version as requirement_version`).
				Joins("left join requirements on schedules.requirement_id = requirements.id").
				Where("schedules.school_id = ?", u.SchoolID).
				Order(clause.OrderBy{Columns: []clause.OrderByColumn{
					{Column: clause.Column{Name: "semester_id"}, Desc: true},
					{Column: clause.Column{Name: "version"}, Desc: true},
				}})

			// 2. Handle standard Go-Admin filtering/sorting
			if semesterID := param.GetFieldValue("semester_id"); semesterID != "" {
				semID, err := strconv.ParseInt(semesterID, 10, 64)
				if err != nil {
					return []map[string]interface{}{
						{"error": err.Error()},
					}, 0
				}
				query = query.Where("semester_id = ?", semID)
			}

			if param.SortField != "" {
				query = query.Order(param.SortField + " " + param.SortType)
			}
			// 3. Handle Global Search (The 'Search' box top-right) and Apply Pagination
			query = query.Offset((param.PageInt - 1) * param.PageSizeInt).Limit(param.PageSizeInt)

			// 4. Execute the query
			if err := query.Find(&reqs).Error; err != nil {
				return []map[string]interface{}{
					{"error": err.Error()},
				}, 0
			}

			// 5. Return the mapped data and the total count for pagination
			result := make([]map[string]interface{}, len(reqs))
			for i := range reqs {
				result[i] = map[string]interface{}{
					"id":                  reqs[i].ID,
					"semester_id":         reqs[i].SemesterID,
					"status":              string(reqs[i].Status),
					"version":             reqs[i].Version.String(),
					"requirement_version": reqs[i].RequirementVersion.String(),
				}
			}
			return result, len(result)
		})

		formList := schedules.GetForm()

		formList.SetTable("schedules").SetTitle("Schedules").SetDescription("Schedules")

		formList.HideResetButton()

		return schedules
	}
}
