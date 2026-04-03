package tables

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/go-playground/validator/v10"
	"github.com/james-wukong/online-school-mgmt/internal/dto"
	model2 "github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	formutils "github.com/james-wukong/online-school-mgmt/internal/utils/form"
	"gorm.io/gorm"
)

func GetTeachersTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		teachers := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		teacherService := services.NewTeacherService(dbConn)
		subService := services.NewSubjectService(dbConn)
		semService := services.NewSemesterService(dbConn)

		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		info := teachers.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("school_id", "=", u.SchoolID)
		}
		info.AddField("Id", "id", db.Int8)
		shoolIDField := info.AddField("School_id", "school_id", db.Int8)
		if !user.IsSuperAdmin() {
			shoolIDField.FieldHide()
		}
		info.AddField("Employee_id", "employee_id", db.Int8)
		info.AddField("Is_active", "is_active", db.Bool).
			FieldEditAble(table2.Switch).
			FieldEditOptions(types.FieldOptions{
				{Value: "true", Text: "Y"}, // 放在第一个代表 on
				{Value: "false", Text: "N"},
			})
		info.AddField("First_name", "first_name", db.Varchar)
		info.AddField("Last_name", "last_name", db.Varchar)
		info.AddField("Email", "email", db.Varchar)
		info.AddField("Phone", "phone", db.Varchar)
		info.AddField("Employment_type", "employment_type", db.Varchar)
		info.AddField("Max_classes_per_day", "max_classes_per_day", db.Int4)
		info.AddField("Hire_date", "hire_date", db.Date)
		info.AddField("Created_at", "created_at", db.Timestamptz)
		info.AddField("Updated_at", "updated_at", db.Timestamptz)
		// Buttons
		info.AddButton(ctx, "Bulk Teachers Create", icon.Tv,
			action.PopUpWithIframe(
				"/teacher/bulk/iframe",
				"Iframe Teachers",
				action.IframeData{
					Src: "/admin/info/bulkteachers/new",
				},
				"900px",
				"600px",
			))

		info.SetTable("teachers").SetTitle("Teachers").SetDescription("Teachers")

		formList := teachers.GetForm()

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

		formList.AddField("Employee_id", "employee_id", db.Int8, form.Text).FieldMust().
			FieldDivider("Teacher Settings")
		formList.AddField("Is_active", "is_active", db.Bool, form.Switch).
			FieldOptions(types.FieldOptions{
				{Text: "Yes", Value: "true"},
				{Text: "No", Value: "false"},
			}).
			FieldDefault("false").
			FieldMust()
		formList.AddField("Subjects", "subjects", db.Varchar, form.SelectBox).
			FieldOptionInitFn(func(val types.FieldModel) types.FieldOptions {
				var c types.FieldOptions
				var teacher *model2.Teachers
				subs, err := subService.List(ctx.Request.Context(), u.SchoolID, 0)
				if err != nil {
					panic(err)
				}
				if val.IsUpdate() {
					tID, err := strconv.ParseInt(val.ID, 10, 64)
					if err != nil {
						panic(err)
					}
					teacher, err = teacherService.GetTeacher(ctx.Request.Context(), tID)
					if err != nil {
						return nil
					}
				}

				for _, s := range subs {
					opt := types.FieldOption{
						Text:  s.Name,
						Value: fmt.Sprint(s.ID)}
					if val.IsUpdate() {
						if exists := slices.ContainsFunc(teacher.Subjects,
							func(t *model2.Subjects) bool {
								return t.ID == s.ID
							}); exists {
							opt.Selected = true
						}
					}
					c = append(c, opt)
				}

				return c
			})

		formList.AddField("First_name", "first_name", db.Varchar, form.Text).FieldMust()
		formList.AddField("Last_name", "last_name", db.Varchar, form.Text).FieldMust()
		formList.AddField("Email", "email", db.Varchar, form.Email)
		formList.AddField("Phone", "phone", db.Varchar, form.Text)
		formList.AddField("Employment_type", "employment_type", db.Varchar, form.SelectSingle).
			// 单选的选项，text代表显示内容，value代表对应值
			FieldOptions(types.FieldOptions{
				{Text: "Permanent", Value: string(model2.Permanent)},
				{Text: "Contract", Value: string(model2.Contract)},
				{Text: "FullTime", Value: string(model2.FullTime)},
				{Text: "PartTime", Value: string(model2.PartTime)},
			}).
			FieldDefault("Permanent")
		formList.AddField("Hire_date", "hire_date", db.Date, form.Date).
			FieldDefault(time.Now().Format(model2.TimeDateLayout)).
			FieldMust()
		formList.AddField("Max_classes_per_day", "max_classes_per_day", db.Int4, form.Number).
			FieldDefault("5")

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
				"/admin/ajax/teacher/sem_timeslot",
				map[string]any{"school_id": fmt.Sprint(u.SchoolID)},
			)).
			FieldMust().
			FieldDivider("Semester Timeslot Settings")
		formList.AddField("Timeslots", "timeslots", db.Varchar, form.SelectBox)

		formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenInsert()
		formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
			FieldHide().FieldNowWhenUpdate()

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// Map values to TeacherRequest struct and validate
			req, err := formutils.MapAndValidate[dto.TeacherCreateRequest](values)
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

			teacher, err := req.ToModel()

			if err != nil {
				return err
			}
			var ts []*model2.TeacherSubjects
			if len(req.SubjectIDs) > 0 {
				for _, subID := range req.SubjectIDs {
					ts = append(ts, &model2.TeacherSubjects{
						TeacherID: teacher.ID,
						SubjectID: subID,
					})
				}
			}
			var tt []*model2.TeacherTimeslots
			if len(req.TimeslotIDs) > 0 {
				for _, slotID := range req.TimeslotIDs {
					tt = append(tt, &model2.TeacherTimeslots{
						TeacherID:  teacher.ID,
						TimeslotID: slotID,
					})
				}
			}
			if err := teacherService.CreateWithAssoc(
				ctx.Request.Context(), teacher, ts, tt); err != nil {
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
				req, err := formutils.MapAndValidate[dto.TeacherStatusUpdateRequest](values)
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
				return teacherService.UpdateStatus(ctx.Request.Context(), req.ToModel())
			}

			// 3. Handle Full Update
			var ts []*model2.TeacherSubjects
			var tt []*model2.TeacherTimeslots
			req, err := formutils.MapAndValidate[dto.TeacherUpdateRequest](values)
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

			fmt.Printf("IsZero: %v\n", time.Time(req.HireDate).IsZero())
			fmt.Printf("hir_date: %s\n", values.Get("hire_date"))
			fmt.Printf("req: %+v\n", req)
			teacher, err := req.ToModel()
			fmt.Printf("teacher: %+v\n", teacher.HireDate)
			if err != nil {
				return err
			}
			if len(req.SubjectIDs) > 0 {
				for _, subID := range req.SubjectIDs {
					ts = append(ts, &model2.TeacherSubjects{
						TeacherID: teacher.ID,
						SubjectID: subID,
					})
				}
			}
			if len(req.TimeslotIDs) > 0 {
				for _, slotID := range req.TimeslotIDs {
					tt = append(tt, &model2.TeacherTimeslots{
						TeacherID:  teacher.ID,
						TimeslotID: slotID,
					})
				}
			}
			//
			var semID int64
			if values.Get("semester_id") != "" {
				semID, err = strconv.ParseInt(values.Get("semester_id"), 10, 64)
				if err != nil {
					return err
				}
			}
			if err := teacherService.UpdateWithAssoc(
				ctx.Request.Context(), teacher, ts, tt, semID); err != nil {
				return err
			}
			return nil
		})

		formList.HideResetButton()
		formList.SetTable("teachers").SetTitle("Teachers").SetDescription("Teachers")

		return teachers
	}
}
