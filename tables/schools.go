package tables

import (
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// School represents the schools table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type School struct {
    // id is GENERATED ALWAYS. We use <-:false to prevent GORM from 
    // including it in INSERT or UPDATE statements.
    ID              int64     `gorm:"column:id;primaryKey;<-:false" json:"id"`
    
    Name            string    `gorm:"column:name;not null;unique" json:"name"`
    Code            string    `gorm:"column:code;not null;unique" json:"code"`
    Address         string    `gorm:"column:address" json:"address"`
    City            string    `gorm:"column:city" json:"city"`
    State           string    `gorm:"column:state" json:"state"`
    PostalCode      string    `gorm:"column:postal_code" json:"postal_code"`
    Country         string    `gorm:"column:country" json:"country"`
    Phone           string    `gorm:"column:phone" json:"phone"`
    Email           string    `gorm:"column:email" json:"email"`
    Website         string    `gorm:"column:website" json:"website"`
    EstablishedYear int       `gorm:"column:established_year" json:"established_year"`
    IsActive        bool      `gorm:"column:is_active;default:true" json:"is_active"`
    
    // CreatedAt is set by the database DEFAULT. We allow read, but restrict write to 'create' only.
    CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    
    // UpdatedAt is handled automatically by GORM's autoUpdateTime feature.
    UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}


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

	formList.SetTable("schools").SetTitle("Schools").SetDescription("Schools")


	return schools
}

