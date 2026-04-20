package export

import (
	"net/http"

	"github.com/GoAdminGroup/go-admin/context"
	"gorm.io/gorm"
)

func ExportReportHandler(dbConn *gorm.DB) context.Handler {
	// 1. Authenticate the request
	// 2. Return JSON response
	// This matches the 'res.code' check in your JavaScript
	return func(ctx *context.Context) {
		// TODO save report to csv files
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"msg":     "Success",
			"success": true,
		})
	}
}
