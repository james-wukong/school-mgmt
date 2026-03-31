package tables

import (
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetTimeslotsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		timeslots := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)

		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := timeslots.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("timeslots.school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8).
			FieldWidth(100)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Sem id", "semester_id", db.Int8).
			FieldWidth(100)
		// 增加字段名 Semester Year
		info.AddField("Sem Year", "year", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_year"]; !ok || r == 0 {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_year"]))
			}).
			FieldSortable().
			FieldWidth(100)
		// 增加字段名 Semester Term
		info.AddField("Sem Term", "semester", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_semester"]; !ok || r == 0 {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				switch value.Row["semesters_goadmin_join_semester"] {
				case "1":
					return template.HTML("Spring")
				case "2":
					return template.HTML("Summer")
				case "3":
					return template.HTML("Fall")
				case "4":
					return template.HTML("Winter")
				default:
					return template.HTML("-")
				}
			}).
			FieldWidth(100)
		// 增加字段名 Semester StartDate
		info.AddField("Sem Start", "start_date", db.Int8).FieldJoin(types.Join{
			Table:     "semesters",   // The table to join with
			Field:     "semester_id", // The foreign key in current table
			JoinField: "id",          // The primary key in joined table
		}).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				if r, ok := value.Row["semesters_goadmin_join_start_date"]; !ok || r == "" {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				// 2. Safely return the value
				return template.HTML(fmt.Sprint(value.Row["semesters_goadmin_join_start_date"]))
			}).
			FieldSortable().
			FieldWidth(100)

		info.AddField("Day_of_week", "day_of_week", db.Int4).
			FieldDisplay(func(value types.FieldModel) interface{} {
				// 1. Check if the joined value is nil
				r, ok := value.Row["day_of_week"]
				if !ok || r == 0 {
					return template.HTML("-") // Return a string or empty template.HTML
				}
				fmt.Printf("r is %#v\n", r)
				switch r {
				case int64(1):
					return template.HTML("Monday")
				case int64(2):
					return template.HTML("Tuesday")
				case int64(3):
					return template.HTML("Wednesday")
				case int64(4):
					return template.HTML("Thursday")
				case int64(5):
					return template.HTML("Friday")
				default:
					return template.HTML("-")
				}
			}).
			FieldWidth(100)
		info.AddField("Start_time", "start_time", db.Time).
			FieldDisplay(func(value types.FieldModel) interface{} {
				t, err := time.Parse(time.RFC3339, value.Value)
				if err != nil {
					panic(err)
				}
				return template.HTML(t.Format(models.TimeSlotLayout))
			}).
			FieldWidth(150)
		info.AddField("End_time", "end_time", db.Time).
			FieldDisplay(func(value types.FieldModel) interface{} {
				t, err := time.Parse(time.RFC3339, value.Value)
				if err != nil {
					panic(err)
				}
				return template.HTML(t.Format(models.TimeSlotLayout))
			})

		// Buttons
		info.AddButton(ctx, "Bulk Timeslots Create", icon.Tv,
			action.PopUpWithIframe(
				"/timeslot/bulk/iframe",
				"Iframe Timeslot",
				action.IframeData{
					Src: "/admin/info/bulktimeslots/new",
				},
				"900px",
				"600px",
			))

		info.SetTable("timeslots").SetTitle("Timeslots").SetDescription("Timeslots")

		formList := timeslots.GetForm()

		semService := services.NewSemesterService(dbConn)

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
			FieldMust().
			FieldDivider("Class Settings")
		formList.AddField("Day_of_week", "day_of_week", db.Int4, form.Radio).
			FieldOptions(types.FieldOptions{
				{Text: "Monday", Value: "1"},
				{Text: "Tuesday", Value: "2"},
				{Text: "Wednesday", Value: "3"},
				{Text: "Thursday", Value: "4"},
				{Text: "Friday", Value: "5"},
				{Text: "Saturday", Value: "6"},
				{Text: "Sunday", Value: "7"},
			}).
			// 设置默认值
			FieldDefault("Monday").
			FieldMust()
		formList.AddField("Start_time", "start_time", db.Time, form.Datetime).
			FieldHelpMsg("修改时分即可").
			FieldMust().
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				t, err := time.Parse(time.RFC3339, value.Value.First())
				if err != nil {
					panic(err)
				}
				return t.Format(models.TimeSlotLayout)
			})
		formList.AddField("End_time", "end_time", db.Time, form.Datetime).
			FieldHelpMsg("修改时分即可").
			FieldMust().
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				t, err := time.Parse(time.RFC3339, value.Value.First())
				if err != nil {
					panic(err)
				}
				return t.Format(models.TimeSlotLayout)
			})

		formList.HideResetButton()
		formList.SetTable("timeslots").SetTitle("Timeslots").SetDescription("Timeslots")

		return timeslots
	}
}
