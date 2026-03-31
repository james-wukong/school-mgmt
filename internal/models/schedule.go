package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// ScheduleStatus handles the Postgres ENUM type-safely
type ScheduleStatus string

const (
	StatusDraft     ScheduleStatus = "Draft"
	StatusPublished ScheduleStatus = "Published"
	StatusActive    ScheduleStatus = "Active"
	StatusArchived  ScheduleStatus = "Archived"
)

type Schedules struct {
	ID            int64           `gorm:"primaryKey;column:id;default:nextval('schedules_id_seq');<-:false" json:"id"`
	SchoolID      int64           `gorm:"column:school_id;not null;index:idx_schedules_school;uniqueIndex:idx_sch_ver_room_time;uniqueIndex:idx_sch_ver_req_time" json:"school_id"`
	SemesterID    int64           `gorm:"column:semester_id;not null;index:idx_requirements_semester" json:"semester_id"`
	RequirementID int64           `gorm:"column:requirement_id;not null;index:idx_schedules_requirement;uniqueIndex:idx_sch_ver_req_time" json:"requirement_id"`
	RoomID        int64           `gorm:"column:room_id;not null;index:idx_schedules_room;uniqueIndex:idx_sch_ver_room_time" json:"room_id"`
	TimeslotID    int64           `gorm:"column:timeslot_id;not null;index:idx_schedules_timeslot;uniqueIndex:idx_sch_ver_room_time;uniqueIndex:idx_sch_ver_req_time" json:"timeslot_id"`
	Status        ScheduleStatus  `gorm:"column:status;type:schedule_status_enum;default:'Draft';index:idx_schedules_status" json:"status"`
	Version       decimal.Decimal `gorm:"column:version;type:numeric(10,2);default:1.00;uniqueIndex:idx_sch_ver_room_time;uniqueIndex:idx_sch_ver_req_time" json:"version"`
	CreatedAt     time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationship mapping
	Requirement Requirements `gorm:"foreignKey:RequirementID" json:"requirement,omitempty"`
	Room        Rooms        `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Semester    Semesters    `gorm:"foreignKey:SemesterID" json:"semester,omitempty"`
	Timeslot    Timeslots    `gorm:"foreignKey:TimeslotID" json:"timeslot,omitempty"`
}
