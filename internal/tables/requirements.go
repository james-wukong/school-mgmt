package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetRequirementsTable(ctx *context.Context) table.Table {
	requirements := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := requirements.GetInfo()

	info.AddField("Id", "id", db.Int8)
	info.AddField("Subject_id", "subject_id", db.Int8)
	info.AddField("Teacher_id", "teacher_id", db.Int8)
	info.AddField("Class_id", "class_id", db.Int8)
	info.AddField("Weekly_sessions", "weekly_sessions", db.Int4)
	info.AddField("Min_day_gap", "min_day_gap", db.Int4)
	info.AddField("Preferred_days", "preferred_days", db.Varchar)

	info.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

	formList := requirements.GetForm()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("Subject_id", "subject_id", db.Int8, form.Text)
	formList.AddField("Teacher_id", "teacher_id", db.Int8, form.Text)
	formList.AddField("Class_id", "class_id", db.Int8, form.Text)
	formList.AddField("Weekly_sessions", "weekly_sessions", db.Int4, form.Number)
	formList.AddField("Min_day_gap", "min_day_gap", db.Int4, form.Number)
	formList.AddField("Preferred_days", "preferred_days", db.Varchar, form.Text)

	formList.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

	return requirements
}
