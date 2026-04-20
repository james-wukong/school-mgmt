package tables

import (
	"fmt"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	kafkautil "github.com/james-wukong/online-school-mgmt/internal/utils/kafka"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetGenerateSchedulesTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		requirements := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

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
		info := requirements.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8).FieldSortable()
		if user.IsSuperAdmin() {
			info.AddField("School_id", "school_id", db.Int8)
		} else {
			info.AddField("School_id", "school_id", db.Int8).FieldHide()
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

		info.AddField("Version", "version", db.Varchar)

		info.AddActionIconButton(ctx, icon.Gears,
			action.Ajax("/admin/schedules/generate",
				func(ctx *context.Context) (success bool, msg string, data interface{}) {
					// 获取参数
					// ctx.FormValue["id"]  选择的id
					// ctx.FormValue["ids"] 选择的id列表，是逗号分割的字符串
					// 1. Get requirement info
					id, err := strconv.ParseInt(ctx.FormValue("id"), 10, 64)
					if err != nil {
						return false, "failure", "id conversion: " + err.Error()
					}
					reqService := services.NewRequirementService(dbConn)
					row, err := reqService.GetByID(ctx.Request.Context(), id)
					if err != nil {
						return false, "failure", "requirement retrieve: " + err.Error()
					}
					// 2. Request to kafka
					if err := kafkautil.ProduceScheduleTask(
						u.SchoolID, row.SemesterID, row.Version.InexactFloat64(), true,
					); err != nil {
						return false, "failure", "kafka producer: " + err.Error()
					}

					return true, "processing", ""
				}),
		)

		info.SetPageSizeList([]int{20, 40, 80, 120}).SetDefaultPageSize(40)
		info.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")
		info.SetGetDataFn(func(param parameter.Parameters) (data []map[string]interface{}, size int) {
			var reqs []*models.RequirementVersion
			// 1. Create the initial query
			query := dbConn.WithContext(ctx.Request.Context()).
				Table("requirements").
				Select("DISTINCT ON (semester_id, version) id, semester_id, version").
				Where("school_id = ?", u.SchoolID).
				Order(clause.OrderBy{Columns: []clause.OrderByColumn{
					{Column: clause.Column{Name: "semester_id"}, Desc: false},
					{Column: clause.Column{Name: "version"}, Desc: false},
				}})
			// 2. Handle standard Go-Admin filtering/sorting
			if param.SortField != "" {
				query = query.Order(param.SortField + " " + param.SortType)
			}

			// 3. Execute the query
			if err := query.Find(&reqs).Error; err != nil {
				return []map[string]interface{}{
					{"error": err.Error()},
				}, 0
			}

			// 4. Return the mapped data and the total count for pagination
			result := make([]map[string]interface{}, len(reqs))
			for i := range reqs {
				result[i] = map[string]interface{}{
					"id":          reqs[i].ID,
					"semester_id": reqs[i].SemesterID,
					"version":     reqs[i].Version.String(),
				}
			}
			return result, len(result)
		})

		return requirements
	}
}
