package tables

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
	// table2 "github.com/GoAdminGroup/go-admin/template/types/table"
)

func GetClassesTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		classes := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		semService := services.NewSemesterService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := classes.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("classes.school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}

		info.AddField("Sem id", "semester_id", db.Int8).FieldWidth(100)

		// 增加字段名 Semester Year
		info.AddField("Sem Year", "year", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_year"]; !ok || r == 0 {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_year"]))
			}).
			FieldSortable().FieldWidth(100)
		// 增加字段名 Semester Term
		info.AddField("Sem Term", "semester", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_semester"]; !ok || r == 0 {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				switch value.Row["semesters_goadmin_join_semester"] {
				case "1":
					return template.HTML("Spring")
				case "2":
					return template.HTML("Summer")
				case "3":
					return template.HTML("Fall")
				case "4":
					return template.HTML("Winter")
				default:
					return template.HTML("-")
				}
			}).FieldWidth(100)
		// 增加字段名 Semester StartDate
		info.AddField("Sem Start", "start_date", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_start_date"]; !ok || r == "" {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_start_date"]))
			}).
			FieldSortable().FieldWidth(100)
		// 增加字段名 Semester EndDate
		info.AddField("Sem End", "end_date", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_end_date"]; !ok || r == "" {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// // 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_end_date"]))
			}).FieldWidth(100)
		info.AddField("Grade", "grade", db.Int4).FieldWidth(100)
		info.AddField("Class", "class", db.Varchar).FieldWidth(100)
		info.AddField("Student_count", "student_count", db.Int4)

		// Buttons
		info.AddButton(ctx, "Bulk Create", icon.Tv,
			action.PopUpWithIframe(
				"/class/bulk/iframe",
				"Iframe Class",
				action.IframeData{
					Src: "/admin/info/bulkclasses/new",
				},
				"900px",
				"600px",
			))
		// info.AddButton(ctx, "ajax", icon.Android, action.Ajax("/admin/ajax",
		// 	func(ctx *context.Context) (success bool, msg string, data interface{}) {
		// 		return true, "请求成功，奥利给", ""
		// 	}))

		info.SetTable("classes").SetTitle("Classes").SetDescription("Classes")

		formList := classes.GetForm()
		formList.AddField("Semester_id", "semester_id", db.Int8, form.SelectSingle).
			FieldOptionInitFn(func(val types.FieldModel) types.FieldOptions {
				var c types.FieldOptions
				s, err := semService.List(ctx.Request.Context(), u.SchoolID, 6)
				if err != nil || len(s) == 0 {
					return nil
				}
				for _, v := range s {
					opt := types.FieldOption{
						Text: fmt.Sprintf(
							"ID: %d, Year: %d, Semester: %d",
							v.ID, v.Year, v.Semester,
						),
						Value: fmt.Sprint(v.ID)}
					if v.ID == val.Row["semester_id"] {
						opt.Selected = true
					}
					c = append(c, opt)
				}

				return c
			}).
			FieldMust().
			FieldDivider("Semester Settings")

		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate().
			FieldMust().
			FieldDivider("Class Settings")
		formList.AddField("Grade", "grade", db.Int4, form.Number).
			FieldMust()
		formList.AddField("Class", "class", db.Varchar, form.Text).
			FieldMust()
		formList.AddField("Student_count", "student_count", db.Int4, form.Number).
			FieldMust()

		formList.HideResetButton()
		formList.HideBackButton()
		formList.SetTable("classes").SetTitle("Classes").SetDescription("Classes")

		return classes
	}
}
