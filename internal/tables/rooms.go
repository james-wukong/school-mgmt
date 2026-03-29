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
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/go-playground/validator/v10"
	"github.com/james-wukong/online-school-mgmt/internal/dto"
	model2 "github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetRoomsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		rooms := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		roomService := services.NewRoomService(dbConn)
		semService := services.NewSemesterService(dbConn)
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
			FieldDivider("Room Settings").
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
		formList.AddField("Semester_id", "semester_id", db.Int8, form.SelectSingle).
			FieldOptionInitFn(func(_ types.FieldModel) types.FieldOptions {
				var c types.FieldOptions
				s, err := semService.List(ctx.Request.Context(), u.SchoolID, 6)
				if err != nil || len(s) == 0 {
					return nil
				}
				for _, v := range s {
					opt := types.FieldOption{
						Text: fmt.Sprintf(
							"ID: %d, Year: %d, Semester: %d",
							v.ID, v.Year, v.Semester,
						),
						Value: fmt.Sprint(v.ID),
					}
					c = append(c, opt)
				}

				return c
			}).
			FieldOnChooseCustom(printDualListBoxJS(
				"semester_id",
				"timeslots[]",
				"/admin/ajax/room/sem_timeslot",
				map[string]any{"school_id": fmt.Sprint(u.SchoolID)},
			)).
			FieldMust().
			FieldDivider("Semester Timeslot Settings")
		formList.AddField("Timeslots", "timeslots", db.Varchar, form.SelectBox)

		formList.AddField("Floor_number", "floor_number", db.Int4, form.Number)
		formList.AddField("Building", "building", db.Varchar, form.Text)
		formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenInsert()
		formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenUpdate()

		formList.HideResetButton()
		formList.SetTable("rooms").SetTitle("Rooms").SetDescription("Rooms")

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// Map values to RoomRequest struct and validate
			req, err := MapAndValidate[dto.RoomCreateRequest](values)
			if err != nil {
				// Check if it's a validation error specifically
				if ve, ok := err.(validator.ValidationErrors); ok {
					// Return the first error found as a string for the UI
					return fmt.Errorf("field '%s' failed validation: %s",
						ve[0].Field(), ve[0].Tag(),
					)
				}
				return err
			}

			room, err := req.ToModel()
			if err != nil {
				return err
			}
			var rt []*model2.RoomTimeslots
			for _, slotID := range req.TimeslotIDs {
				rt = append(rt, &model2.RoomTimeslots{
					RoomID:     room.ID,
					TimeslotID: slotID,
				})
			}
			if err := roomService.CreateWithAssoc(
				ctx.Request.Context(), room, rt); err != nil {
				return err
			}

			return nil
		})

		// 取代更新函数
		formList.SetUpdateFn(func(values form2.Values) error {
			// 1. Identify the Record
			id := values.Get("id") // Ensure 'id' is in your Info/Form fields
			if id == "" {
				return errors.New("missing primary key for update")
			}

			// 2. Handle Single Update (The Switch Toggle)
			if len(values) == 2 && values.Has("is_active") {
				req, err := MapAndValidate[dto.RoomStatusUpdateRequest](values)
				if err != nil {
					// Check if it's a validation error specifically
					if ve, ok := err.(validator.ValidationErrors); ok {
						// Return the first error found as a string for the UI
						return fmt.Errorf("field '%s' failed validation: %s",
							ve[0].Field(), ve[0].Tag(),
						)
					}
					return err
				}
				return roomService.UpdateStatus(ctx.Request.Context(), req.ToModel())
			}

			// 3. Handle Full Update
			var rt []*model2.RoomTimeslots
			req, err := MapAndValidate[dto.RoomUpdateRequest](values)
			if err != nil {
				// Check if it's a validation error specifically
				if ve, ok := err.(validator.ValidationErrors); ok {
					// Return the first error found as a string for the UI
					return fmt.Errorf("field '%s' failed validation: %s",
						ve[0].Field(), ve[0].Tag(),
					)
				}
				return err
			}
			room, err := req.ToModel()
			if err != nil {
				return err
			}
			for _, slotID := range req.TimeslotIDs {
				rt = append(rt, &model2.RoomTimeslots{
					RoomID:     room.ID,
					TimeslotID: slotID,
				})
			}
			//
			semID, err := strconv.ParseInt(values.Get("semester_id"), 10, 64)
			if err != nil {
				return err
			}
			if err := roomService.UpdateWithAssoc(
				ctx.Request.Context(), room, rt, semID); err != nil {
				return err
			}
			return nil
		})

		return rooms
	}
}
