package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
)

type ClassBase struct {
	// School context is usually required for every class
	SemesterID   int64  `form:"semester_id" csv:"semester_id" json:"semester_id" validate:"required,gt=0"`
	SchoolID     int64  `form:"school_id" csv:"school_id" json:"school_id" validate:"required,gt=0"`
	Grade        int    `form:"grade" csv:"grade" json:"grade" validate:"required,gt=0"`
	ClassName    string `form:"class" csv:"class" json:"class" validate:"required,min=2,max=100"`
	StudentCount int    `form:"student_count" csv:"student_count" json:"student_count" validate:"omitempty,gt=0"`
}

type ClassCreateRequest struct {
	ClassBase
}

type ClassUpdateRequest struct {
	ID int64 `form:"id" validate:"required"` // The ID is mandatory
	ClassBase
}

func (req *ClassBase) toModel() (*models.Classes, error) {
	return &models.Classes{
		SchoolID:     req.SchoolID,
		SemesterID:   req.SemesterID,
		Grade:        req.Grade,
		ClassName:    req.ClassName,
		StudentCount: req.StudentCount,
	}, nil
}

func (req *ClassCreateRequest) ToModel() (*models.Classes, error) {
	return req.toModel()
}

func (req *ClassUpdateRequest) ToModel() (*models.Classes, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}
