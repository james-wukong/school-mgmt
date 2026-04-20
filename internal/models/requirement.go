package models

import (
	"github.com/shopspring/decimal"
)

type Requirements struct {
	ID             int64           `gorm:"primaryKey;column:id;default:nextval('requirements_id_seq');<-:false" json:"id"`
	SchoolID       int64           `gorm:"column:school_id;not null;index:idx_requirements_school" json:"school_id"`
	SemesterID     int64           `gorm:"column:semester_id;not null;index:idx_requirements_semester" json:"semester_id"`
	SubjectID      int64           `gorm:"column:subject_id;not null;index:idx_requirements_subject" json:"subject_id"`
	TeacherID      int64           `gorm:"column:teacher_id;not null;index:idx_requirements_teacher" json:"teacher_id"`
	ClassID        int64           `gorm:"column:class_id;not null;index:idx_requirements_class" json:"class_id"`
	WeeklySessions int             `gorm:"column:weekly_sessions;not null;default:1" json:"weekly_sessions"`
	MinDayGap      int             `gorm:"column:min_day_gap;not null;default:0" json:"min_day_gap"`
	PreferredDays  *string         `gorm:"column:preferred_days;type:varchar(100)" json:"preferred_days"` // Pointer handles NULL
	Version        decimal.Decimal `gorm:"column:version;type:numeric(10,2);default:1.00" json:"version"`

	// Relationships for Eager Loading (Preload)
	School   *Schools   `gorm:"foreignKey:SchoolID" json:"school,omitempty"`
	Semester *Semesters `gorm:"foreignKey:SemesterID" json:"semester,omitempty"`
	Subject  *Subjects  `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher  *Teachers  `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Class    *Classes   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
}

type RequirementVersion struct {
	ID         int64           `gorm:"primaryKey;column:id;default:nextval('requirements_id_seq')" json:"id"`
	SemesterID int64           `gorm:"column:semester_id;not null;index:idx_requirements_semester" json:"semester_id"`
	Version    decimal.Decimal `gorm:"column:version;type:numeric(10,2);default:1.00" json:"version"`
}
