package tables

import (
	"fmt"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"
)

func GetRequirementsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		requirements := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		semService := services.NewSemesterService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}
		semesters, err := semService.List(ctx.Request.Context(), u.SchoolID, 6)
		if err != nil {
			panic(err)
		}
		var semFilterOpts types.FieldOptions
		for _, semester := range semesters {
			semFilterOpts = append(semFilterOpts, types.FieldOption{
				Value: fmt.Sprint(semester.ID),
				Text: fmt.Sprintf("ID: %d - Year: %d - Semester: %d",
					semester.ID, semester.Year, semester.Semester,
				),
			})
		}
		info := requirements.GetInfo()
		if !user.IsSuperAdmin() {
			info = info.Where("requirements.school_id", "=", u.SchoolID)
		}

		info.AddField("Id", "id", db.Int8).FieldSortable()
		if user.IsSuperAdmin() {
			info.AddField("School_id", "school_id", db.Int8)
		} else {
			info.AddField("School_id", "school_id", db.Int8).FieldHide()
		}

		info.AddField("Version", "version", db.Varchar).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorEqual,
				Width:    450,
			}).
			FieldSortable()
		info.AddField("Semester", "semester_id", db.Int8).
			FieldDisplay(func(model types.FieldModel) interface{} {
				semID, err := strconv.ParseInt(model.Value, 10, 64)
				if err != nil {
					panic(err)
				}
				sem, err := semService.GetByID(ctx.Request.Context(), semID)
				if err != nil {
					panic(err)
				}
				return fmt.Sprintf("ID: %s - Year: %d - Semester: %d",
					model.Value, sem.Year, sem.Semester,
				)
			}).
			FieldFilterable(types.FilterType{
				FormType: form.SelectSingle,
				Width:    450,
			}).
			FieldFilterOptions(semFilterOpts).
			FieldSortable()
		// TODO add subject info
		info.AddField("Subject_id", "subject_id", db.Int8).FieldSortable()
		info.AddField("SubjectName", "name", db.Varchar).
			// Force the field to reference the JOINED table column
			FieldJoin(types.Join{
				BaseTable: "requirements", // The main table
				Field:     "subject_id",   // The field in the main table
				Table:     "subjects",     // The join table
				JoinField: "id",           // The ID in the join table
			}).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorLike,
				Width:    450,
			})

		info.AddField("SubjectCode", "code", db.Varchar).
			// Force the field to reference the JOINED table column
			FieldJoin(types.Join{
				BaseTable: "requirements", // The main table
				Field:     "subject_id",   // The field in the main table
				Table:     "subjects",     // The join table
				JoinField: "id",           // The ID in the join table
			}).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorLike,
				Width:    450,
			}).
			FieldSortable()
		// TODO add teacher info
		info.AddField("Teacher_id", "teacher_id", db.Int8).FieldSortable()
		info.AddField("Teacher FirstName", "first_name", db.Varchar).
			// Force the field to reference the JOINED table column
			FieldJoin(types.Join{
				BaseTable: "requirements", // The main table
				Field:     "teacher_id",   // The field in the main table
				Table:     "teachers",     // The join table
				JoinField: "id",           // The ID in the join table
			}).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorLike,
				Width:    450,
			}).
			FieldSortable()
		info.AddField("Teacher LastName", "last_name", db.Varchar).
			// Force the field to reference the JOINED table column
			FieldJoin(types.Join{
				BaseTable: "requirements", // The main table
				Field:     "teacher_id",   // The field in the main table
				Table:     "teachers",     // The join table
				JoinField: "id",           // The ID in the join table
			}).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorLike,
				Width:    450,
			}).
			FieldSortable()
		// TODO add class info
		info.AddField("Class_id", "class_id", db.Int8).
			FieldSortable()
		info.AddField("ClassGrade", "grade", db.Varchar).
			// Force the field to reference the JOINED table column
			FieldJoin(types.Join{
				BaseTable: "requirements", // The main table
				Field:     "class_id",     // The field in the main table
				Table:     "classes",      // The join table
				JoinField: "id",           // The ID in the join table
			}).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorEqual,
				Width:    450,
			})
		info.AddField("ClassName", "class", db.Varchar).
			// Force the field to reference the JOINED table column
			FieldJoin(types.Join{
				BaseTable: "requirements", // The main table
				Field:     "class_id",     // The field in the main table
				Table:     "classes",      // The join table
				JoinField: "id",           // The ID in the join table
			}).
			FieldFilterable(types.FilterType{
				FormType: form.Text,
				Operator: types.FilterOperatorLike,
				Width:    450,
			})
		info.AddField("Weekly_sessions", "weekly_sessions", db.Int4)
		info.AddField("Min_day_gap", "min_day_gap", db.Int4)
		info.AddField("Preferred_days", "preferred_days", db.Varchar)
		// Buttons
		info.AddButton(ctx, "批量创建", icon.Tv,
			action.PopUpWithIframe(
				"/requirement/bulk/iframe",
				"Iframe Requirement",
				action.IframeData{
					Src: "/admin/info/bulkrequirements/new",
				},
				"900px",
				"600px",
			))

		info.SetPageSizeList([]int{20, 40, 80, 120}).SetDefaultPageSize(40)
		info.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

		formList := requirements.GetForm()
		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate()
		schoolField := formList.AddField("School_id", "school_id", db.Int8, form.Default).
			FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
				if value.IsCreate() {
					return fmt.Sprint(u.SchoolID)
				}
				return value.Value.First()
			})
		// Apply the conditional visibility
		if !user.IsSuperAdmin() {
			schoolField.FieldHide()
		}
		formList.AddField("Semester_id", "semester_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Subject_id", "subject_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Teacher_id", "teacher_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Class_id", "class_id", db.Int8, form.Text).FieldMust()
		formList.AddField("Weekly_sessions", "weekly_sessions", db.Int4, form.Number).
			FieldDefault("5").
			FieldMust()
		formList.AddField("Min_day_gap", "min_day_gap", db.Int4, form.Number).
			FieldDefault("0")
		formList.AddField("Preferred_days", "preferred_days", db.Varchar, form.Text)

		formList.HideResetButton()
		formList.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

		formList.SetPreProcessFn(func(values form2.Values) form2.Values {
			for k, v := range values {
				fmt.Printf("k is %s and values is %+v\n", k, v)
			}
			return values
		})

		return requirements
	}
}
