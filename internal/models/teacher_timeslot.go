package models

type TeacherTimeslots struct {
	// Foreign Key to Teachers
	TeacherID int64 `gorm:"not null;primaryKey" json:"teacher_id"`
	// Foreign Key to Timeslots
	TimeslotID int64 `gorm:"not null;primaryKey" json:"timeslot_id"`
}
