package tables

import (
	"path/filepath"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"gorm.io/gorm"
)

func GetSampleFilesTable(dbConn *gorm.DB) table.Generator {
	return func(ctx *context.Context) table.Table {
		sampleFiles := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

		info := sampleFiles.GetInfo()

		info.AddField("Id", "id", db.Int8)

		info.AddField("Filename", "filename", db.Varchar).
			FieldDisplay(func(value types.FieldModel) interface{} {
				return filepath.Base(value.Value)
			}).
			FieldDownLoadable("http://127.0.0.1:8091/uploads/samples/")
		info.AddField("Description", "description", db.Text)

		info.HideNewButton()
		info.HideEditButton()
		info.HideExportButton()
		info.HideDeleteButton()
		info.HideDetailButton()

		info.SetTable("sample_files").SetTitle("Samplefiles").SetDescription("Samplefiles")

		return sampleFiles
	}
}
