package tables

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/james-wukong/online-school-mgmt/internal/dto"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/provider"
	"github.com/james-wukong/online-school-mgmt/internal/services"
	formutils "github.com/james-wukong/online-school-mgmt/internal/utils/form"
	"gorm.io/gorm"

	// table2 "github.com/GoAdminGroup/go-admin/template/types/table"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
)

func GetBulkSubjectsTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		subjects := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))
		user := auth.Auth(ctx)
		userService := services.NewAdminUserService(dbConn)
		u, err := userService.GetUserSchoolID(ctx.Request.Context(), user.Id)
		if err != nil {
			panic(err)
		}

		formList := subjects.GetForm()

		formList.AddField("Choose Source Type", "source_type", db.Tinyint, form.SelectSingle).
			FieldOptions(types.FieldOptions{
				{Text: "JSON", Value: "0"},
				{Text: "CSV", Value: "1"},
			}).
			FieldOnChooseHide("1", "json").
			FieldOnChooseShow("0", "json").
			FieldOnChooseShow("1", "csv").
			FieldDefault("1").
			FieldDivider("Source Settings")

		formList.AddField("JSON", "json", db.Int, form.TextArea).
			FieldDefault(printSampleSubjectJSON()).
			FieldHelpMsg(`采用json格式: 参考默认值`)
		formList.AddField("CSV", "csv", db.Int, form.File).
			FieldOptionExt(map[string]interface{}{
				"allowClear": true,
			})

		formList.AddField("Id", "id", db.Int8, form.Default).
			FieldDisableWhenCreate().
			FieldHide()
		formList.AddField("School Id", "school_id", db.Int8, form.Default).
			FieldDefault(fmt.Sprint(u.SchoolID)).
			FieldHide()

		// 取代新增函数
		formList.SetInsertFn(func(values form2.Values) error {
			// values 为传入的表单参数
			// 1. Get the request object from the context
			// values.Context is the *context.Context provided by GoAdmin
			var reader provider.DataReader[dto.SubjectCreateRequest]
			var filePath string

			// when user chooses JSON
			switch values.Get("source_type") {
			case "0":
				text := values.Get("json")
				if text == "" {
					return errors.New("empty json")
				}
				reader = provider.NewTextReader[dto.SubjectCreateRequest](text)

			// when user chooses File
			case "1":
				filePath := "./uploads/" + values.Get("csv")
				if filePath == "" {
					return errors.New("no file uploaded")
				}
				csvValidation := provider.FileValidationConfig{
					AllowedExtensions: []string{".csv"},
					MaxSizeBytes:      5 * 1024 * 1024, // 5 MB
				}
				// Validate before processing
				if err := provider.ValidateFile(filePath, csvValidation); err != nil {
					return err // GoAdmin shows this message to the user
				}

				reader = provider.NewCSVReader[dto.SubjectCreateRequest](filePath, false)

			default:
				return errors.New("unsupported source type")
			}

			// Always clean up the uploaded file when done
			defer func() {
				if values.Get("csv") != "" {
					if err := os.Remove("./uploads/" + values.Get("csv")); err != nil {
						log.Printf("cleanup failed for %s: %v", filePath, err)
					}
				}
			}()

			// Start DB Process
			var subs []*models.Subjects
			var errs []error

			rows, err := reader.Read(ctx.Request.Context())
			if err != nil {
				return err
			}
			subService := services.NewSubjectService(dbConn)

			schID, err := strconv.ParseInt(values.Get("school_id"), 10, 64)
			if err != nil {
				return err
			}

			for i, row := range rows {
				// Validate struct
				if err := formutils.Validate.Struct(row); err != nil {
					errs = append(errs, fmt.Errorf("row %d error: %v", i, err))
				}
				// reconstruct requirement struct
				row.SchoolID = schID

				r, err := row.ToModel()
				if err != nil {
					return err
				}
				subs = append(subs, r)
			}
			if len(errs) > 0 {
				// TODO log errors or email errors
				for _, e := range errs {
					fmt.Printf("err: %s\n", e.Error())
				}
				return errors.New("check console for detail errors")
			}

			return subService.FilterAndCreateInBatches(ctx.Request.Context(), subs)
		})

		formList.HideContinueNewCheckBox()
		formList.HideResetButton()
		formList.SetTable("requirements").SetTitle("Requirements").SetDescription("Requirements")

		return subjects
	}
}
