package models

import (
	"time"
)

const (
	TimeSlotLayout = "15:04"
)

type DayOfWeek int

const (
	Monday    DayOfWeek = iota + 1 // 1
	Tuesday                        // 2
	Wednesday                      // 3
	Thursday                       // 4
	Friday                         // 5
	Saturday
	Sunday
)

// Timeslots represents the timeslot table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Timeslots struct {
	// id is GENERATED ALWAYS. use <-:false to prevent GORM from
	// including it in INSERT or UPDATE statements.
	ID int64 `gorm:"column:id;primaryKey;<-:false" json:"id"`

	// Foreign Key to School
	SchoolID   int64 `gorm:"column:school_id;not null;uniqueIndex:idx_rooms_school_code" json:"school_id"`
	SemesterID int64 `gorm:"not null;index:idx_timeslots_semester;index:idx_timeslot_unique,unique" json:"semester_id"`
	// DayOfWeek: 1->Monday, 2->Tuesday, etc.
	DayOfWeek DayOfWeek `gorm:"not null;index:idx_timeslots_day;index:idx_timeslot_unique,unique" json:"day_of_week"`
	// StartDate and EndDate use time.Time.
	StartTime time.Time `gorm:"type:time;not null;index:idx_timeslot_unique,unique" json:"start_time"`
	EndTime   time.Time `gorm:"type:time;not null" json:"end_time"`

	// Relationships (Optional but recommended for Eager Loading)
	School   *Schools   `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`
	Semester *Semesters `gorm:"foreignKey:SemesterID;constraint:OnDelete:CASCADE;" json:"semester,omitempty"`
}

func NewTimeslots(schoolID, semesterID int64,
	day DayOfWeek,
	startTime time.Time,
) *Timeslots {
	return &Timeslots{
		SchoolID:   schoolID,
		SemesterID: semesterID,
		DayOfWeek:  day,
		StartTime:  startTime,
	}
}
