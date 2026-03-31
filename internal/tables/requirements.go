package tables

import (
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
)

func GetRequirementsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		requirements := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		// roomService := services.NewRoomService(dbConn)
		// semService := services.NewSemesterService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := requirements.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8)
		info.AddField("School_id", "school_id", db.Int8)
		info.AddField("Semester_id", "semester_id", db.Int8)
		info.AddField("Subject_id", "subject_id", db.Int8)
		info.AddField("Teacher_id", "teacher_id", db.Int8)
		info.AddField("Class_id", "class_id", db.Int8)
		info.AddField("Weekly_sessions", "weekly_sessions", db.Int4)
		info.AddField("Min_day_gap", "min_day_gap", db.Int4)
		info.AddField("Preferred_days", "preferred_days", db.Varchar)
		info.AddField("Version", "version", db.Varchar)
		// Buttons
		info.AddButton(ctx, "Bulk Requirements Create", icon.Tv,
			action.PopUpWithIframe(
				"/requirement/bulk/iframe",
				"Iframe Requirement",
				action.IframeData{
					Src: "/admin/info/bulkrequirements/new",
				},
				"900px",
				"600px",
			))

		info.SetPageSizeList([]int{20, 40, 80, 120}).SetDefaultPageSize(40)
		info.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

		formList := requirements.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				if value.IsCreate() {
					return u.SchoolID
				}
				return value.Value
			})
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}
		formList.AddField("Semester_id", "semester_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Subject_id", "subject_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Teacher_id", "teacher_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Class_id", "class_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Weekly_sessions", "weekly_sessions", db.Int4, form.Number).
			FieldDefault("5").
			FieldMust()
		formList.AddField("Min_day_gap", "min_day_gap", db.Int4, form.Number).
			FieldDefault("0")
		formList.AddField("Preferred_days", "preferred_days", db.Varchar, form.Text)

		formList.HideResetButton()
		formList.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

		return requirements
	}
}
