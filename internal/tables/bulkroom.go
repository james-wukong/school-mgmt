package tables

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
	// table2 "github.com/GoAdminGroup/go-admin/template/types/table"
)

func GetBulkRoomsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		rooms := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		roomService := services.NewRoomService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}

		formList := rooms.GetForm()

		formList.AddField("Building Name", "building_name", db.Varchar, form.Text).
			FieldMust().
			FieldDivider("Building Settings")
		formList.AddField("Total Floors", "total_floor", db.Int4, form.Number).
			FieldMust().
			FieldDefault("3")
		formList.AddField("Number of Rooms/Floor", "num_of_rooms", db.Int4, form.Number).
			FieldMust().
			FieldDefault("10")
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate().
			FieldHide()
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldDefault(fmt.Sprint(u.SchoolID))
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}

		formList.HideResetButton()
		formList.HideContinueNewCheckBox()
		formList.SetTable("rooms").SetTitle("Rooms").SetDescription("Rooms")

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// values 为传入的表单参数
			// 1. validate input
			if values.IsEmpty("school_id", "building_name", "total_floor", "num_of_rooms") {
				return errors.New("fields are missing")
			}
			// Convert the string to int64
			// base 10 (decimal), bitSize 64 (for int64)
			if err != nil {
				return err
			}
			totalFloor, err := strconv.ParseInt(values.Get("total_floor"), 10, 64)
			if err != nil {
				return err
			}
			numOfRooms, err := strconv.ParseInt(values.Get("num_of_rooms"), 10, 64)
			if err != nil {
				return err
			}

			return roomService.CreateRoomsInBatches(ctx.Request.Context(),
				values.Get("building_name"), u.SchoolID,
				int(totalFloor), int(numOfRooms))
		})
		return rooms
	}
}
