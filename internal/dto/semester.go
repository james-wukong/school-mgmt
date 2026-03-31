package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/types"
)

type SemesterBase struct {
	// School context is usually required for every semester
	SchoolID int64 `form:"school_id" csv:"school_id" json:"school_id" validate:"required,gt=0"`

	Year      int             `form:"year" cvs:"year" json:"year" validate:"required"`
	Semester  int             `form:"semester" cvs:"semester" json:"semester" validate:"required"`
	StartDate types.CivilDate `form:"start_date" csv:"start_date" json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   types.CivilDate `form:"end_date" csv:"end_date" json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type SemesterCreateRequest struct {
	SemesterBase
}

type SemesterUpdateRequest struct {
	ID int64 `form:"id" csv:"id" json:"id" validate:"required"` // The ID is mandatory
	SemesterBase
}

func (req *SemesterBase) toModel() (*models.Semesters, error) {
	return &models.Semesters{
		SchoolID:  req.SchoolID,
		Year:      req.Year,
		Semester:  req.Semester,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func (req *SemesterCreateRequest) ToModel() (*models.Semesters, error) {
	return req.toModel()
}

func (req *SemesterUpdateRequest) ToModel() (*models.Semesters, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}
