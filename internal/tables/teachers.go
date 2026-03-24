package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	model2 "github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetTeachersTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		teachers := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := teachers.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Employee_id", "employee_id", db.Int8)
		info.AddField("Is_active", "is_active", db.Bool).
			FieldEditAble(table2.Switch).
			FieldEditOptions(types.FieldOptions{
				{Value: "true", Text: "Y"}, // 放在第一个代表 on
				{Value: "false", Text: "N"},
			})
		info.AddField("First_name", "first_name", db.Varchar)
		info.AddField("Last_name", "last_name", db.Varchar)
		info.AddField("Email", "email", db.Varchar)
		info.AddField("Phone", "phone", db.Varchar)
		info.AddField("Employment_type", "employment_type", db.Varchar)
		info.AddField("Max_classes_per_day", "max_classes_per_day", db.Int4)
		info.AddField("Hire_date", "hire_date", db.Date)
		info.AddField("Created_at", "created_at", db.Timestamptz)
		info.AddField("Updated_at", "updated_at", db.Timestamptz)

		info.SetTable("teachers").SetTitle("Teachers").SetDescription("Teachers")

		formList := teachers.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				return u.SchoolID
			})
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}

		formList.AddField("Employee_id", "employee_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Is_active", "is_active", db.Bool, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "Yes", Value: "true"},
				{Text: "No", Value: "false"},
			}).
			FieldDefault("false").
			FieldMust()

		formList.AddField("First_name", "first_name", db.Varchar, form.Text).FieldMust()
		formList.AddField("Last_name", "last_name", db.Varchar, form.Text).FieldMust()
		formList.AddField("Email", "email", db.Varchar, form.Email)
		formList.AddField("Phone", "phone", db.Varchar, form.Text)
		formList.AddField("Employment_type", "employment_type", db.Varchar, form.Select).
			// 单选的选项，text代表显示内容，value代表对应值
			FieldOptions(types.FieldOptions{
				{Text: "Permanent", Value: string(model2.Permanent)},
				{Text: "Contract", Value: string(model2.Contract)},
				{Text: "FullTime", Value: string(model2.FullTime)},
				{Text: "PartTime", Value: string(model2.PartTime)},
			}).
			FieldDefault(string(model2.Permanent))
		formList.AddField("Hire_date", "hire_date", db.Date, form.Date)
		formList.AddField("Max_classes_per_day", "max_classes_per_day", db.Int4, form.Number)
		formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenInsert()
		formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenUpdate()

		formList.SetTable("teachers").SetTitle("Teachers").SetDescription("Teachers")

		return teachers
	}
}
