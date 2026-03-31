package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetSchoolsTable(ctx *context.Context) table.Table {
	schools := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := schools.GetInfo().SetPrimaryKey("id", db.Bigint)

	info.AddField("Id", "id", db.Int8)
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Code", "code", db.Varchar)
	info.AddField("Established_year", "established_year", db.Int4)
	info.AddField("Is_active", "is_active", db.Bool).FieldBool("true", "false")
	info.AddField("Phone", "phone", db.Varchar)
	info.AddField("Email", "email", db.Varchar)
	info.AddField("State", "state", db.Varchar)
	info.AddField("Postal_code", "postal_code", db.Varchar)
	info.AddField("Country", "country", db.Varchar)
	info.AddField("Website", "website", db.Varchar)
	info.AddField("Address", "address", db.Text)
	info.AddField("City", "city", db.Varchar)
	info.AddField("Created_at", "created_at", db.Timestamptz)
	info.AddField("Updated_at", "updated_at", db.Timestamptz)

	info.SetTable("schools").SetTitle("Schools").SetDescription("Schools")

	formList := schools.GetForm()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate().
		FieldNotAllowEdit()
	formList.AddField("Name", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Code", "code", db.Varchar, form.Text).FieldMust()
	formList.AddField("State", "state", db.Varchar, form.Text)
	formList.AddField("Postal_code", "postal_code", db.Varchar, form.Text)
	formList.AddField("Country", "country", db.Varchar, form.Text)
	formList.AddField("Phone", "phone", db.Varchar, form.Text)
	formList.AddField("Email", "email", db.Varchar, form.Email)
	formList.AddField("Website", "website", db.Varchar, form.Text)
	formList.AddField("Address", "address", db.Text, form.Text)
	formList.AddField("City", "city", db.Varchar, form.Text)
	formList.AddField("Established_year", "established_year", db.Int4, form.Number)
	formList.AddField("Is_active", "is_active", db.Bool, form.Switch).
		FieldOptions(types.FieldOptions{
			{Text: "Active", Value: "true"},
			{Text: "InActive", Value: "false"},
		}).
		FieldDefault("false").
		FieldMust()
	formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenUpdate()

	formList.HideResetButton()
	formList.SetTable("schools").SetTitle("Schools").SetDescription("Schools")

	return schools
}
