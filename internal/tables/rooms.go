package tables

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	model2 "github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetRoomsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		rooms := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		slotService := services.NewTimeslotService(dbConn)
		roomService := services.NewRoomService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := rooms.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8)

		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Name", "name", db.Varchar)
		info.AddField("Code", "code", db.Varchar)
		info.AddField("Room_type", "room_type", db.Varchar)
		info.AddField("Capacity", "capacity", db.Int4)
		info.AddField("Is_active", "is_active", db.Bool).
			FieldEditAble(table2.Switch).
			FieldEditOptions(types.FieldOptions{
				{Value: "true", Text: "Y"}, // 放在第一个代表 on
				{Value: "false", Text: "N"},
			})
		info.AddField("Building", "building", db.Varchar)
		info.AddField("Floor_number", "floor_number", db.Int4)
		info.AddField("Created_at", "created_at", db.Timestamptz)
		info.AddField("Updated_at", "updated_at", db.Timestamptz)

		info.SetTable("rooms").SetTitle("Rooms").SetDescription("Rooms")

		formList := rooms.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				return fmt.Sprint(u.SchoolID)
			})
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}
		formList.AddField("Code", "code", db.Varchar, form.Text).
			FieldHelpMsg("必填, 唯一").
			FieldMust()
		formList.AddField("Name", "name", db.Varchar, form.Text).
			FieldHelpMsg("必填, 唯一").
			FieldMust()
		formList.AddField("Room_type", "room_type", db.Varchar, form.SelectSingle).
			// 单选的选项，text代表显示内容，value代表对应值
			FieldOptions(types.FieldOptions{
				{Text: "Regular", Value: string(model2.Regular)},
				{Text: "Lab", Value: string(model2.Lab)},
				{Text: "Gym", Value: string(model2.Gym)},
			}).
			FieldDefault("Regular").
			FieldMust()
		formList.AddField("Capacity", "capacity", db.Int4, form.Number).
			FieldDefault("50").
			FieldMust()
		formList.AddField("Is_active", "is_active", db.Bool, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "Yes", Value: "true"},
				{Text: "No", Value: "false"},
			}).
			FieldDefault("false").
			FieldMust()
		formList.AddField("Timeslots", "timeslots", db.Varchar, form.SelectBox).
			FieldOptionInitFn(func(val types.FieldModel) types.FieldOptions {
				var c types.FieldOptions
				var room *model2.Rooms
				slots, err := slotService.List(ctx.Request.Context(), u.SchoolID, 0, 0)
				if err != nil {
					panic(err)
				}
				if val.IsUpdate() {
					rID, err := strconv.ParseInt(val.ID, 10, 64)
					if err != nil {
						panic(err)
					}
					room, err = roomService.GetRoom(ctx.Request.Context(), rID)
					if err != nil {
						return nil
					}
				}

				for _, s := range slots {
					opt := types.FieldOption{
						Text: fmt.Sprintf("%d: %s-%s",
							s.DayOfWeek,
							s.StartTime.Format(model2.TimeSlotLayout),
							s.EndTime.Format(model2.TimeSlotLayout),
						),
						Value: fmt.Sprint(s.ID)}
					if val.IsUpdate() {
						if exists := slices.ContainsFunc(room.Timeslots,
							func(t *model2.Timeslots) bool {
								return t.ID == s.ID
							}); exists {
							opt.Selected = true
						}
					}
					c = append(c, opt)
				}

				return c
			})
		formList.AddField("Floor_number", "floor_number", db.Int4, form.Number)
		formList.AddField("Building", "building", db.Varchar, form.Text)
		formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenInsert()
		formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenUpdate()

		formList.HideResetButton()
		formList.HideBackButton()
		formList.SetTable("rooms").SetTitle("Rooms").SetDescription("Rooms")

		return rooms
	}
}
