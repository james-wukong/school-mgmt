package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/shopspring/decimal"
)

type RequirementBase struct {
	SchoolID       int64 `form:"school_id" csv:"school_id" json:"school_id" validate:"required"`
	SemesterID     int64 `form:"semester_id" csv:"semester_id" json:"semester_id" validate:"required"`
	SubjectID      int64 `form:"subject_id" csv:"subject_id" json:"subject_id" validate:"required"`
	TeacherID      int64 `form:"teacher_id" csv:"teacher_id" json:"teacher_id" validate:"required"`
	ClassID        int64 `form:"class_id" csv:"class_id" json:"class_id" validate:"required"`
	WeeklySessions int   `form:"weekly_sessions" csv:"weekly_sessions" json:"weekly_sessions" validate:"required,min=1"`
	MinDayGap      int   `form:"min_day_gap" csv:"min_day_gap" json:"min_day_gap" validate:"min=0"`

	// PreferredDays is a comma-separated string from the frontend (e.g., "1,2,3")
	// We use *string to allow it to be NULL in the database if empty.
	PreferredDays *string `form:"preferred_days" csv:"preferred_days" json:"preferred_days" validate:"omitempty"`

	// Version ensures we are editing the correct iteration of the requirements.
	Version decimal.Decimal `form:"version" csv:"version" json:"version" validate:"required"`

	// relationships
	Semester SemesterCreateRequest `csv:"semester_assoc_,inline" json:"-" validate:"semester,omitempty"`
	Subject  SubjectCreateRequest  `csv:"subject_assoc_,inline" json:"-" validate:"subject,omitempty"`
	Teacher  TeacherCreateRequest  `csv:"teacher_assoc_,inline" json:"-" validate:"teacher,omitempty"`
	Class    ClassCreateRequest    `csv:"class_assoc_,inline" json:"-" validate:"class,omitempty"`
}

type RequirementCreateRequest struct {
	RequirementBase
}

type RequirementUpdateRequest struct {
	ID int64 `form:"id" csv:"id" json:"id" validate:"required"` // The ID is mandatory
	RequirementBase
}

func (req *RequirementBase) toModel() (*models.Requirements, error) {
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
	return &models.Requirements{
		SchoolID:       req.SchoolID,
		SemesterID:     req.SemesterID,
		SubjectID:      req.SubjectID,
		TeacherID:      req.TeacherID,
		ClassID:        req.ClassID,
		WeeklySessions: req.WeeklySessions,
		MinDayGap:      req.MinDayGap,
		PreferredDays:  req.PreferredDays,
		Version:        req.Version,
		Semester:       sem,
		Subject:        sub,
		Teacher:        tch,
		Class:          cls,
	}, nil
}

func (req *RequirementCreateRequest) ToModel() (*models.Requirements, error) {
	return req.toModel()
}

func (req *RequirementUpdateRequest) ToModel() (*models.Requirements, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}
