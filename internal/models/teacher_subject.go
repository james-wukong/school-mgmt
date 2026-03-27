package models

type TeacherSubjects struct {
	// Foreign Key to Teachers
	TeacherID int64 `gorm:"not null;primaryKey" json:"teacher_id"`
	// Foreign Key to Subjects
	SubjectID int64 `gorm:"not null;primaryKey" json:"subject_id"`
}
