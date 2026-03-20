package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetSemestersTable(ctx *context.Context) table.Table {

	semesters := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := semesters.GetInfo()

	info.AddField("Id", "id", db.Int8)
	info.AddField("School_id", "school_id", db.Int8)
	info.AddField("Year", "year", db.Int4)
	info.AddField("Semester", "semester", db.Int4)
	info.AddField("Start_date", "start_date", db.Date)
	info.AddField("End_date", "end_date", db.Date)

	info.SetTable("semesters").SetTitle("Semesters").SetDescription("Semesters")

	formList := semesters.GetForm()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("School_id", "school_id", db.Int8, form.Text)
	formList.AddField("Year", "year", db.Int4, form.Number)
	formList.AddField("Semester", "semester", db.Int4, form.Number)
	formList.AddField("Start_date", "start_date", db.Date, form.Datetime)
	formList.AddField("End_date", "end_date", db.Date, form.Datetime)

	formList.SetTable("semesters").SetTitle("Semesters").SetDescription("Semesters")

	return semesters
}
