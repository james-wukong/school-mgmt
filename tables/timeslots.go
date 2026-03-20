package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetTimeslotsTable(ctx *context.Context) table.Table {

	timeslots := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := timeslots.GetInfo()

	info.AddField("Id", "id", db.Int8)
	info.AddField("Semester_id", "semester_id", db.Int8)
	info.AddField("Day_of_week", "day_of_week", db.Int4)
	info.AddField("Start_time", "start_time", db.Time)
	info.AddField("End_time", "end_time", db.Time)

	info.SetTable("timeslots").SetTitle("Timeslots").SetDescription("Timeslots")

	formList := timeslots.GetForm()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("Semester_id", "semester_id", db.Int8, form.Text)
	formList.AddField("Day_of_week", "day_of_week", db.Int4, form.Number)
	formList.AddField("Start_time", "start_time", db.Time, form.Datetime)
	formList.AddField("End_time", "end_time", db.Time, form.Datetime)

	formList.SetTable("timeslots").SetTitle("Timeslots").SetDescription("Timeslots")

	return timeslots
}
