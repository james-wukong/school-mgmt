package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetStudentsTable(ctx *context.Context) table.Table {
	students := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := students.GetInfo()

	info.AddField("Updated_at", "updated_at", db.Timestamptz)
	info.AddField("School_id", "school_id", db.Int8)
	info.AddField("Date_of_birth", "date_of_birth", db.Date)
	info.AddField("Admission_date", "admission_date", db.Date)
	info.AddField("Is_active", "is_active", db.Bool)
	info.AddField("Created_at", "created_at", db.Timestamptz)
	info.AddField("Id", "id", db.Int8)
	info.AddField("Blood_group", "blood_group", db.Varchar)
	info.AddField("Gender", "gender", db.Varchar)
	info.AddField("Student_number", "student_number", db.Varchar)
	info.AddField("First_name", "first_name", db.Varchar)
	info.AddField("Last_name", "last_name", db.Varchar)
	info.AddField("Email", "email", db.Varchar)
	info.AddField("Phone", "phone", db.Varchar)

	info.SetTable("students").SetTitle("Students").SetDescription("Students")

	formList := students.GetForm()
	formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenUpdate()
	formList.AddField("School_id", "school_id", db.Int8, form.Text)
	formList.AddField("Date_of_birth", "date_of_birth", db.Date, form.Datetime)
	formList.AddField("Admission_date", "admission_date", db.Date, form.Datetime)
	formList.AddField("Is_active", "is_active", db.Bool, form.Text)
	formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("Blood_group", "blood_group", db.Varchar, form.Text)
	formList.AddField("Gender", "gender", db.Varchar, form.Text)
	formList.AddField("Student_number", "student_number", db.Varchar, form.Text)
	formList.AddField("First_name", "first_name", db.Varchar, form.Text)
	formList.AddField("Last_name", "last_name", db.Varchar, form.Text)
	formList.AddField("Email", "email", db.Varchar, form.Email)
	formList.AddField("Phone", "phone", db.Varchar, form.Text)

	formList.SetTable("students").SetTitle("Students").SetDescription("Students")

	return students
}
