package tables

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
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
		sch, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := classes.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", sch.ID)
		}

		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}

		info.AddField("Semester_id", "semester_id", db.Int8)

		// 增加字段名 Semester Year
		info.AddField("Semester_Year", "year", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_year"]; !ok || r == 0 {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// // 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_year"]))
			})
		// 增加字段名 Semester Term
		info.AddField("Semester_Term", "semester", db.Int8).FieldJoin(types.Join{
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
			})
		// 增加字段名 Semester StartDate
		info.AddField("Semester_Term", "start_date", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_start_date"]; !ok || r == "" {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// // 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_start_date"]))
			})
		// 增加字段名 Semester EndDate
		info.AddField("Semester_Term", "end_date", db.Int8).FieldJoin(types.Join{
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
			})
		info.AddField("Grade", "grade", db.Int4)
		info.AddField("Class", "class", db.Varchar)
		info.AddField("Student_count", "student_count", db.Int4)

		info.SetTable("classes").SetTitle("Classes").SetDescription("Classes")

		formList := classes.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
		formList.AddField("Semester_id", "semester_id", db.Int8, form.Default)
		formList.AddField("Semester Year", "year", db.Int4, form.Number).FieldDefault("2000")
		formList.AddField("Semester Year", "semester", db.Int4, form.Number).FieldDefault("1")
		formList.AddField("Grade", "grade", db.Int4, form.Number)
		formList.AddField("Student_count", "student_count", db.Int4, form.Number)
		formList.AddField("Class", "class", db.Varchar, form.Text)

		formList.SetTable("classes").SetTitle("Classes").SetDescription("Classes")

		return classes
	}
}
