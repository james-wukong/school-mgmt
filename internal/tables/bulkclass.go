package tables

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
		formList.AddField("Grade-Class", "grades", db.Int4, form.Text).
			FieldMust().
			FieldHelpMsg("eg: 1-3, 2-5,3-8, 左边是年级，右边是班级数量，年级和班级用-分割，年级之间用，分割").
			FieldDivider("Class Settings")

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// values 为传入的表单参数
			var cls []*model2.Classes
			// 1. validate input
			if values.IsEmpty("semester_id", "grades") {
				return errors.New("semester id and grades can not be empty")
			}
			// Convert the string to int64
			// base 10 (decimal), bitSize 64 (for int64)
			sID, err := strconv.ParseInt(values.Get("semester_id"), 10, 64)

			if err != nil {
				return err
			}
			sem, err := semService.GetByID(ctx.Request.Context(), sID)
			if err != nil {
				return err
			}
			reg := regexp.MustCompile("[^0-9-,]+")

			// Replace all characters matching the regex with an empty string
			pairs := reg.ReplaceAllString(values.Get("grades"), "")
			for _, v := range strings.Split(pairs, ",") {
				if v == "" {
					continue
				}
				pair := strings.Split(v, "-")
				g, err := strconv.ParseInt(pair[0], 10, 0)
				if err != nil {
					return err
				}
				total, err := strconv.ParseInt(pair[1], 10, 0)
				if err != nil {
					return err
				}
				for t := range total {
					cls = append(cls, &model2.Classes{
						SemesterID: sID,
						Grade:      int(g),
						ClassName:  strconv.Itoa(int(t) + 1),
					})
				}
			}
			err = semService.AppendClasses(ctx.Request.Context(), sem, cls)
			if err != nil {
				return err
			}
			return nil
		})

		formList.HideResetButton()
		formList.SetTable("classes").SetTitle("Classes").SetDescription("Classes")
		return classes
	}
}
