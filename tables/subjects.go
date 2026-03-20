package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetSubjectsTable(ctx *context.Context) table.Table {

	subjects := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := subjects.GetInfo()

	info.AddField("Required_lab", "required_lab", db.Bool)
	info.AddField("School_id", "school_id", db.Int8)
	info.AddField("Is_heavy", "is_heavy", db.Bool)
	info.AddField("Id", "id", db.Int8)
	info.AddField("Code", "code", db.Varchar)
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Description", "description", db.Text)

	info.SetTable("subjects").SetTitle("Subjects").SetDescription("Subjects")

	formList := subjects.GetForm()
	formList.AddField("Required_lab", "required_lab", db.Bool, form.Text)
	formList.AddField("School_id", "school_id", db.Int8, form.Text)
	formList.AddField("Is_heavy", "is_heavy", db.Bool, form.Text)
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("Code", "code", db.Varchar, form.Text)
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Description", "description", db.Text, form.RichText)

	formList.SetTable("subjects").SetTitle("Subjects").SetDescription("Subjects")

	return subjects
}
