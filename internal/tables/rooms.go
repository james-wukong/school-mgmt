package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetRoomsTable(ctx *context.Context) table.Table {

	rooms := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := rooms.GetInfo()

	info.AddField("Updated_at", "updated_at", db.Timestamptz)
	info.AddField("School_id", "school_id", db.Int8)
	info.AddField("Capacity", "capacity", db.Int4)
	info.AddField("Floor_number", "floor_number", db.Int4)
	info.AddField("Is_active", "is_active", db.Bool)
	info.AddField("Created_at", "created_at", db.Timestamptz)
	info.AddField("Id", "id", db.Int8)
	info.AddField("Code", "code", db.Varchar)
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Room_type", "room_type", db.Varchar)
	info.AddField("Building", "building", db.Varchar)
	info.AddField("Available_days", "available_days", db.Varchar)

	info.SetTable("rooms").SetTitle("Rooms").SetDescription("Rooms")

	formList := rooms.GetForm()
	formList.AddField("Updated_at", "updated_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenUpdate()
	formList.AddField("School_id", "school_id", db.Int8, form.Text)
	formList.AddField("Capacity", "capacity", db.Int4, form.Number)
	formList.AddField("Floor_number", "floor_number", db.Int4, form.Number)
	formList.AddField("Is_active", "is_active", db.Bool, form.Text)
	formList.AddField("Created_at", "created_at", db.Timestamptz, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField("Id", "id", db.Int8, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("Code", "code", db.Varchar, form.Text)
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Room_type", "room_type", db.Varchar, form.Text)
	formList.AddField("Building", "building", db.Varchar, form.Text)
	formList.AddField("Available_days", "available_days", db.Varchar, form.Text)

	formList.SetTable("rooms").SetTitle("Rooms").SetDescription("Rooms")

	return rooms
}
