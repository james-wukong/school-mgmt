package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetSchedulesTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		schedules := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := schedules.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}
		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Requirement ID", "requirement_id", db.Int8)
		info.AddField("Status", "status", db.Varchar)
		info.AddField("Version", "version", db.Varchar)

		info.SetTable("schedules").SetTitle("Schedules").SetDescription("Schedules")

		formList := schedules.GetForm()

		formList.SetTable("schedules").SetTitle("Schedules").SetDescription("Schedules")

		formList.HideResetButton()

		return schedules
	}
}
