package middlewares

import (
	"net/http"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

func UserContextInjector(dbConn db.Connection) context.Handler {
	return func(ctx *context.Context) {
	user := auth.Auth(ctx)

	// Not logged in
	if user.Id == 0 {
		// Set a response so the user knows WHY it was aborted
		ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code": http.StatusUnauthorized,
			"message":  "You do not have permission to access this region.",
		})
		ctx.Abort()
		return
	}

	res, err := db.WithDriver(dbConn).
						Table("goadmin_users").
						Where("id", "=", user.Id).
						Select("school_id").
						First()
	if err == nil {
		// ✅ inject custom field
		ctx.SetUserValue("schoolID", res["school_id"])
	}

	ctx.Next()
}
}