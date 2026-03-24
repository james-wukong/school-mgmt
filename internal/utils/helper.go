package utils

import (
	stdctx "context"

	adminctx "github.com/GoAdminGroup/go-admin/context"
)

func GetStdContext(ctx *adminctx.Context) stdctx.Context {
	return ctx.Request.Context()
}
