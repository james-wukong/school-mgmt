package tables

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	model2 "github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	"gorm.io/gorm"

	// table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
)

func GetBulkClassesTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		classes := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		semService := services.NewSemesterService(dbConn)
		clsService := services.NewClassService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}

		formList := classes.GetForm()
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
					c = append(c, opt)
				}
				if len(c) > 0 {
					c[0].Selected = true
				}

				return c
			}).
			FieldMust().
			FieldDivider("Semester Settings")

		formList.AddTable("New", "grade_class", func(panel *types.FormPanel) {
			panel.AddField("Grade", "grade", db.Int, form.Number).
				FieldDefault("1").
				FieldHideLabel()
			panel.AddField("Number of Classes", "classes", db.Varchar, form.Number).
				FieldDefault("8").
				FieldHideLabel()
		})

		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate().
			FieldHide()

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// values 为传入的表单参数
			var cls []*model2.Classes

			for k, v := range values {
				fmt.Printf("k is %s and value is: %+v", k, v)
			}
			// 1. validate input
			if values.IsEmpty("semester_id") {
				return errors.New("semester id and grades can not be empty")
			}
			if len(values["grade"]) > 0 && len(values["grade"]) != len(values["classes"]) {
				return errors.New("number of grade and classes don't match")
			}
			// Convert the string to int64
			// base 10 (decimal), bitSize 64 (for int64)
			sID, err := strconv.ParseInt(values.Get("semester_id"), 10, 64)
			if err != nil {
				return err
			}
			for i, g := range values["grade"] {
				ig, err := strconv.ParseInt(g, 10, 64)
				if err != nil {
					return err
				}
				numOfClass, err := strconv.ParseInt(values["classes"][i], 10, 64)
				if err != nil {
					return err
				}
				for c := range numOfClass {
					class := &model2.Classes{
						SemesterID:   sID,
						Grade:        int(ig),
						ClassName:    fmt.Sprintf("%02d", c+1),
						SchoolID:     u.SchoolID,
						StudentCount: 40,
					}
					fmt.Printf("class: %+v\n", class)
					cls = append(cls, class)
				}
			}

			return clsService.CreateInBatches(ctx.Request.Context(), cls)
		})

		formList.HideResetButton()
		formList.HideContinueNewCheckBox()
		formList.SetTable("classes").SetTitle("Classes").SetDescription("Classes")
		return classes
	}
}
