package tables

import (
	"fmt"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetSemestersTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		semesters := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := semesters.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Year", "year", db.Int4)
		info.AddField("Semester", "semester", db.Int4).
			FieldDisplay(func(model types.FieldModel) interface{} {
				switch {
				case model.Value == "1":
					return "Spring"
				case model.Value == "2":
					return "Summer"
				case model.Value == "3":
					return "Fall"
				case model.Value == "4":
					return "Winter"
				default:
					return "-"
				}
			})
		info.AddField("Start_date", "start_date", db.Date)
		info.AddField("End_date", "end_date", db.Date)

		info.SetTable("semesters").SetTitle("Semesters").SetDescription("Semesters")
		formList := semesters.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				return u.SchoolID
			}).FieldMust()
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}
		formList.AddField("Year", "year", db.Int4, form.Number).
			FieldDefault(fmt.Sprint(time.Now().Year())).
			FieldMust()
		formList.AddField("Semester", "semester", db.Int4, form.Radio).
			FieldOptions(types.FieldOptions{
				{Text: "Spring", Value: "1"},
				{Text: "Summer", Value: "2"},
				{Text: "Fall", Value: "3"},
				{Text: "Winter", Value: "4"},
			}).
			// 设置默认值
			FieldDefault("Fall").
			FieldMust()
		formList.AddField("Start_date", "start_date", db.Date, form.Date).
			FieldMust()
		formList.AddField("End_date", "end_date", db.Date, form.Date).
			FieldMust()

		formList.SetTable("semesters").SetTitle("Semesters").SetDescription("Semesters")

		return semesters
	}
}
