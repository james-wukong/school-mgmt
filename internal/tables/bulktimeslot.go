package tables

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/dto"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"

	// table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
)

func GetBulkTimeslotsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		timeslots := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		semService := services.NewSemesterService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}

		formList := timeslots.GetForm()
		formList.AddField("Semester_id", "semester_id", db.Int8, form.SelectSingle).
			FieldOptionInitFn(func(val types.FieldModel) types.FieldOptions {
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
						Value: fmt.Sprint(v.ID)}
					if v.ID == val.Row["semester_id"] {
						opt.Selected = true
					}
					c = append(c, opt)
				}

				return c
			}).
			FieldMust().
			FieldDivider("Semester Settings")

		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate().
			FieldHide()
		formList.AddField("Timeslots", "timeslots", db.Int4, form.TextArea).
			FieldDefault(`
{
  "Monday": [
    {"start_time": "09:00", "end_time": "09:45"},
    {"start_time": "10:00", "end_time": "10:45"},
    {"start_time": "11:00", "end_time": "11:45"},
    {"start_time": "13:00", "end_time": "13:45"},
    {"start_time": "14:00", "end_time": "14:45"},
    {"start_time": "15:00", "end_time": "15:45"}
  ],
  "Tuesday": [
    {"start_time": "09:00", "end_time": "09:45"},
    {"start_time": "10:00", "end_time": "10:45"},
    {"start_time": "11:00", "end_time": "11:45"},
    {"start_time": "13:00", "end_time": "13:45"},
    {"start_time": "14:00", "end_time": "14:45"},
    {"start_time": "15:00", "end_time": "15:45"}
  ],
  "Wednesday": [
    {"start_time": "09:00", "end_time": "09:45"},
    {"start_time": "10:00", "end_time": "10:45"},
    {"start_time": "11:00", "end_time": "11:45"},
    {"start_time": "13:00", "end_time": "13:45"},
    {"start_time": "14:00", "end_time": "14:45"},
    {"start_time": "15:00", "end_time": "15:45"}
  ],
  "Thursday": [
    {"start_time": "09:00", "end_time": "09:45"},
    {"start_time": "10:00", "end_time": "10:45"},
    {"start_time": "11:00", "end_time": "11:45"},
    {"start_time": "13:00", "end_time": "13:45"},
    {"start_time": "14:00", "end_time": "14:45"},
    {"start_time": "15:00", "end_time": "15:45"}
  ],
  "Friday": [
    {"start_time": "09:00", "end_time": "09:45"},
    {"start_time": "10:00", "end_time": "10:45"},
    {"start_time": "11:00", "end_time": "11:45"},
    {"start_time": "13:00", "end_time": "13:45"},
    {"start_time": "14:00", "end_time": "14:45"},
    {"start_time": "15:00", "end_time": "15:45"}
  ]
}
`).
			FieldMust().
			FieldHelpMsg(`采用json格式: {"day":[{"start_time": time, "end_time": time}]}`).
			FieldDivider("Timeslot Settings")

		formList.SetTable("timeslots").SetTitle("Timeslots").SetDescription("Timeslots")

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// values 为传入的表单参数
			// 1. validate input
			if values.IsEmpty("semester_id", "timeslots") {
				return errors.New("semester id and timeslots can not be empty")
			}
			// Convert the string to int64
			// base 10 (decimal), bitSize 64 (for int64)
			semID, err := strconv.ParseInt(values.Get("semester_id"), 10, 64)

			if err != nil {
				return err
			}
			sem, err := semService.GetByID(ctx.Request.Context(), semID)
			if err != nil {
				return err
			}
			// Create an instance of the struct
			var timetable dto.Schedule

			// Convert the JSON string to a byte slice and unmarshal into the struct
			err = json.Unmarshal([]byte(values.Get("timeslots")), &timetable)
			if err != nil {
				return err
			}
			// check input data is consistent
			if err := timetable.Validate(); err != nil {
				return err
			}
			// Now we are safe to save data into database
			sem.Timeslots = timetable.MapToTimeslots(semID, sem.SchoolID)
			err = semService.ReplaceWithTimeslotAssoc(ctx.Request.Context(), sem)

			return err
		})

		return timeslots
	}
}
