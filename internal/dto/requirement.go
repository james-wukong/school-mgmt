package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/shopspring/decimal"
)

type RequirementBase struct {
	SchoolID       int64 `form:"school_id" csv:"school_id" json:"school_id" validate:"-"`
	SemesterID     int64 `form:"semester_id" csv:"semester_id" json:"semester_id" validate:"-"`
	SubjectID      int64 `form:"subject_id" csv:"subject_id" json:"subject_id" validate:"-"`
	TeacherID      int64 `form:"teacher_id" csv:"teacher_id" json:"teacher_id" validate:"-"`
	ClassID        int64 `form:"class_id" csv:"class_id" json:"class_id" validate:"-"`
	WeeklySessions *int  `form:"weekly_sessions" csv:"weekly_sessions" json:"weekly_sessions" validate:"required,min=1,max=30"`
	MinDayGap      *int  `form:"min_day_gap" csv:"min_day_gap" json:"min_day_gap" validate:"required,min=0"`

	// PreferredDays is a comma-separated string from the frontend (e.g., "1,2,3")
	// We use *string to allow it to be NULL in the database if empty.
	PreferredDays *string `form:"preferred_days" csv:"preferred_days" json:"preferred_days" validate:"omitempty"`

	// Version ensures we are editing the correct iteration of the requirements.
	Version decimal.Decimal `form:"version" csv:"version" json:"version" validate:"-"`
}

type RequirementCreateRequest struct {
	RequirementBase

	// relationships
	Semester SemesterUpdateRequest `csv:"semester_assoc_,inline" json:"-" validate:"-"`
	Subject  SubjectUpdateRequest  `csv:"subject_assoc_,inline" json:"subject" validate:"required,omitempty"`
	Teacher  TeacherUpdateRequest  `csv:"teacher_assoc_,inline" json:"teacher" validate:"required,omitempty"`
	Class    ClassUpdateRequest    `csv:"class_assoc_,inline" json:"class" validate:"required,omitempty"`
}

type RequirementUpdateRequest struct {
	ID int64 `form:"id" csv:"id" json:"id" validate:"required"` // The ID is mandatory
	RequirementBase

	// relationships
	Semester SemesterUpdateRequest `csv:"semester_assoc_,inline" json:"-" validate:"-"`
	Subject  SubjectUpdateRequest  `csv:"subject_assoc_,inline" json:"subject" validate:"required,omitempty"`
	Teacher  TeacherUpdateRequest  `csv:"teacher_assoc_,inline" json:"teacher" validate:"required,omitempty"`
	Class    ClassUpdateRequest    `csv:"class_assoc_,inline" json:"class" validate:"required,omitempty"`
}

func (req *RequirementBase) toModel() (*models.Requirements, error) {
	m := &models.Requirements{
		SchoolID:      req.SchoolID,
		SemesterID:    req.SemesterID,
		SubjectID:     req.SubjectID,
		TeacherID:     req.TeacherID,
		ClassID:       req.ClassID,
		PreferredDays: req.PreferredDays,
		Version:       req.Version,
	}
	// Safe Check: WeeklySessions
	if req.WeeklySessions != nil {
		m.WeeklySessions = *req.WeeklySessions
	}

	// Safe Check: MinDayGap
	if req.MinDayGap != nil {
		m.MinDayGap = *req.MinDayGap
	}
	return m, nil
}

func (req *RequirementCreateRequest) ToModel() (*models.Requirements, error) {

	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	sem, err := req.Semester.ToModel()
	if err != nil {
		return nil, err
	}
	sub, err := req.Subject.ToModel()
	if err != nil {
		return nil, err
	}
	tch, err := req.Teacher.ToModel()
	if err != nil {

		return nil, err
	}
	cls, err := req.Class.ToModel()
	if err != nil {
		return nil, err
	}
	if sem.ID != 0 {
		m.Semester = sem
	}
	if sub.ID != 0 {
		m.Subject = sub
	}
	if tch.ID != 0 {
		m.Teacher = tch
	}
	if cls.ID != 0 {
		m.Class = cls
	}
	return m, nil
}

func (req *RequirementUpdateRequest) ToModel() (*models.Requirements, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	sem, err := req.Semester.ToModel()
	if err != nil {
		return nil, err
	}
	sub, err := req.Subject.ToModel()
	if err != nil {
		return nil, err
	}
	tch, err := req.Teacher.ToModel()
	if err != nil {
		return nil, err
	}
	cls, err := req.Class.ToModel()
	if err != nil {
		return nil, err
	}
	m.Semester = sem
	m.Subject = sub
	m.Teacher = tch
	m.Class = cls
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID

	return m, nil
}
