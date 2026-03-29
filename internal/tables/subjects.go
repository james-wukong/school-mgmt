package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetSubjectsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		subjects := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := subjects.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Code", "code", db.Varchar)
		info.AddField("Name", "name", db.Varchar)
		info.AddField("Requires_lab", "requires_lab", db.Bool).
			FieldEditAble(table2.Switch).
			FieldEditOptions(types.FieldOptions{
				{Value: "true", Text: "Y"}, // 放在第一个代表 on
				{Value: "false", Text: "N"},
			})
		info.AddField("Is_heavy", "is_heavy", db.Bool).
			FieldEditAble(table2.Switch).
			FieldEditOptions(types.FieldOptions{
				{Value: "true", Text: "Y"}, // 放在第一个代表 on
				{Value: "false", Text: "N"},
			})
		info.AddField("Description", "description", db.Text)

		info.SetTable("subjects").SetTitle("Subjects").SetDescription("Subjects")

		formList := subjects.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
			// Define the field once
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				return u.SchoolID
			})
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}
		formList.AddField("Code", "code", db.Varchar, form.Text).FieldMust()
		formList.AddField("Name", "name", db.Varchar, form.Text).FieldMust()
		formList.AddField("Requires_lab", "requires_lab", db.Bool, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "Yes", Value: "true"},
				{Text: "No", Value: "false"},
			}).
			FieldDefault("false").
			FieldMust()
		formList.AddField("Is_heavy", "is_heavy", db.Bool, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "Yes", Value: "true"},
				{Text: "No", Value: "false"},
			}).
			FieldDefault("false").
			FieldMust()
		formList.AddField("Description", "description", db.Text, form.RichText)

		formList.HideResetButton()
		formList.SetTable("subjects").SetTitle("Subjects").SetDescription("Subjects")

		return subjects
	}
}
