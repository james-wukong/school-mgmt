package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
)

type SubjectBase struct {
	// School context is usually required for every subject
	SchoolID    int64  `form:"school_id" csv:"school_id" json:"school_id" validate:"required,gt=0"`
	Name        string `form:"name" csv:"name" json:"name" validate:"required,min=2,max=100"`
	Code        string `form:"code" csv:"code" json:"code" validate:"required,min=2,max=100"`
	Description string `form:"description" csv:"description" json:"description" validate:"omitempty"`
	RequiresLab bool   `form:"requires_lab" csv:"requires_lab" json:"requires_lab"`
	IsHeavy     bool   `form:"is_heavy" csv:"is_heavy" json:"is_heavy"`
}

type SubjectCreateRequest struct {
	SubjectBase
}

type SubjectUpdateRequest struct {
	ID int64 `form:"id" validate:"required"` // The ID is mandatory
	SubjectBase
}

func (req *SubjectBase) toModel() (*models.Subjects, error) {

	return &models.Subjects{
		SchoolID:    req.SchoolID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		RequiresLab: req.RequiresLab,
		IsHeavy:     req.IsHeavy,
	}, nil
}

func (req *SubjectCreateRequest) ToModel() (*models.Subjects, error) {
	return req.toModel()
}

func (req *SubjectUpdateRequest) ToModel() (*models.Subjects, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}
