package models

import (
	"time"
)

type RoomType string

const (
	Regular RoomType = "Regular"
	Lab     RoomType = "Lab"
	Gym     RoomType = "Gym"
)

// Rooms represents the physical rooms/facilities within a school.
// It uses GORM tags to map structural logic to the database schema.
type Rooms struct {
	ID          int64     `gorm:"primaryKey;autoIncrement:true;autoIncrementIncrement:1;<-:false" json:"id"`
	SchoolID    int64     `gorm:"not null;index:idx_school_code,unique" json:"school_id"`
	Code        string    `gorm:"type:varchar(50);not null;index:idx_school_code,unique" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	RoomType    RoomType  `gorm:"type:varchar(50)" json:"room_type"` // e.g., 'Classroom', 'Lab'
	Capacity    int       `gorm:"default:40" json:"capacity"`
	FloorNumber *int      `json:"floor_number"` // Pointer allows for NULL values
	Building    string    `gorm:"type:varchar(100)" json:"building"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	// BelongTo relationship: A Room belongs to a School
	School *Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;" json:"school,omitempty"`

	// Many-to-Many relationship with Subjects
	// The 'many2many' tag points to the join table name in PostgreSQL
	// foreignKey: Primary Key of "Source" -> Teachers.ID
	// joinForeignKey: Column in Join Table for Source -> teacher_subjects.teacher_id
	// references: Primary Key of "Target" -> Subjects.ID
	// joinReferences: Column in Join Table for Target -> teacher_subjects.subject_id
	Timeslots []*Timeslots `gorm:"many2many:room_timeslots;foreignKey:ID;joinForeignKey:RoomID;References:ID;joinReferences:TimeslotID;constraint:OnDelete:CASCADE;" json:"timeslots,omitempty"`
}
