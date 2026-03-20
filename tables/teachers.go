package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetTeachersTable(ctx *context.Context) table.Table {

	teachers := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := teachers.GetInfo()

	info.AddField("Id", "id", db.Int8)
	info.AddField("School_id", "school_id", db.Int8)
	info.AddField("Employee_id", "employee_id", db.Int8)
	info.AddField("Hire_date", "hire_date", db.Date)
	info.AddField("Max_classes_per_day", "max_classes_per_day", db.Int4)
	info.AddField("Is_active", "is_active", db.Bool)
	info.AddField("Created_at", "created_at", db.Timestamptz)
	info.AddField("Updated_at", "updated_at", db.Timestamptz)
	info.AddField("First_name", "first_name", db.Varchar)
	info.AddField("Last_name", "last_name", db.Varchar)
	info.AddField("Email", "email", db.Varchar)
	info.AddField("Phone", "phone", db.Varchar)
	info.AddField("Employment_type", "employment_type", db.Varchar)

	info.SetTable("teachers").SetTitle("Teachers").SetDescription("Teachers")

	formList := teachers.GetForm()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("School_id", "school_id", db.Int8, form.Text)
	formList.AddField("Employee_id", "employee_id", db.Int8, form.Text)
	formList.AddField("Hire_date", "hire_date", db.Date, form.Datetime)
	formList.AddField("Max_classes_per_day", "max_classes_per_day", db.Int4, form.Number)
	formList.AddField("Is_active", "is_active", db.Bool, form.Text)
	formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenUpdate()
	formList.AddField("First_name", "first_name", db.Varchar, form.Text)
	formList.AddField("Last_name", "last_name", db.Varchar, form.Text)
	formList.AddField("Email", "email", db.Varchar, form.Email)
	formList.AddField("Phone", "phone", db.Varchar, form.Text)
	formList.AddField("Employment_type", "employment_type", db.Varchar, form.Text)

	formList.SetTable("teachers").SetTitle("Teachers").SetDescription("Teachers")

	return teachers
}
