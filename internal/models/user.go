package models

import (
	"time"
)

type AdminUser struct {
	// ID uses a custom sequence in Postgres, so we mark it as primaryKey
	// and let the database handle the default value.
	ID            uint32 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username      string `gorm:"column:username;type:varchar(100);not null;unique" json:"username"`
	Password      string `gorm:"column:password;type:varchar(100);not null" json:"password"`
	Name          string `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Avatar        string `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	RememberToken string `gorm:"column:remember_token;type:varchar(100)" json:"remember_token"`

	// Add your custom fields with GORM tags
	SchoolID int64    `gorm:"column:school_id;not null" json:"school_id"` // Belongs To School
	School   *Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;" json:"school,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:now()" json:"updated_at"`
}

func (AdminUser) TableName() string {
	return "goadmin_users"
}
