package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

func GetSchedulesTable(ctx *context.Context) table.Table {
	schedules := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := schedules.GetInfo().HideFilterArea()

	info.SetTable("schedules").SetTitle("Schedules").SetDescription("Schedules")

	formList := schedules.GetForm()

	formList.SetTable("schedules").SetTitle("Schedules").SetDescription("Schedules")

	return schedules
}
