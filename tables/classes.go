package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetClassesTable(ctx *context.Context) table.Table {

	classes := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := classes.GetInfo()

	info.AddField("Id", "id", db.Int8)
	info.AddField("Semester_id", "semester_id", db.Int8)
	info.AddField("Grade", "grade", db.Int4)
	info.AddField("Student_count", "student_count", db.Int4)
	info.AddField("Class", "class", db.Varchar)

	info.SetTable("classes").SetTitle("Classes").SetDescription("Classes")

	formList := classes.GetForm()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("Semester_id", "semester_id", db.Int8, form.Text)
	formList.AddField("Grade", "grade", db.Int4, form.Number)
	formList.AddField("Student_count", "student_count", db.Int4, form.Number)
	formList.AddField("Class", "class", db.Varchar, form.Text)

	formList.SetTable("classes").SetTitle("Classes").SetDescription("Classes")

	return classes
}
