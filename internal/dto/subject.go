package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
)

type SubjectBase struct {
	// School context is usually required for every subject
	SchoolID    int64  `form:"school_id" csv:"school_id" json:"school_id" validate:"-"`
	Name        string `form:"name" csv:"name" json:"name" validate:"required,min=2,max=100"`
	Description string `form:"description" csv:"description" json:"description" validate:"omitempty"`
	RequiresLab *bool  `form:"requires_lab" csv:"requires_lab" json:"requires_lab"`
	IsHeavy     *bool  `form:"is_heavy" csv:"is_heavy" json:"is_heavy"`
}

type ReqSubjectCreateRequest struct {
	SubjectBase
	Code string `form:"code" csv:"code" json:"code" validate:"omitempty,min=2,max=100"`
}

type ReqSubjectUpdateRequest struct {
	ID   int64  `form:"id" csv:"id" json:"id" validate:"required"` // The ID is mandatory
	Code string `form:"code" csv:"code" json:"code" validate:"omitempty,min=2,max=100"`
	SubjectBase
}

type SubjectCreateRequest struct {
	SubjectBase
	Code string `form:"code" csv:"code" json:"code" validate:"required,min=2,max=100"`
}

type SubjectUpdateRequest struct {
	ID int64 `form:"id" csv:"id" json:"id" validate:"required"` // The ID is mandatory
	SubjectBase
}

func (req *SubjectBase) toModel() (*models.Subjects, error) {
	m := &models.Subjects{
		SchoolID:    req.SchoolID,
		Name:        req.Name,
		Description: req.Description,
	}
	// Safe Check: RequiresLab
	if req.RequiresLab != nil {
		m.RequiresLab = *req.RequiresLab
	}

	// Safe Check: IsHeavy
	if req.IsHeavy != nil {
		m.IsHeavy = *req.IsHeavy
	}
	return m, nil
}

func (req *SubjectCreateRequest) ToModel() (*models.Subjects, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}

	m.Code = req.Code

	return m, nil
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

func (req *ReqSubjectUpdateRequest) ToModel() (*models.Subjects, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.Code = req.Code
	m.ID = req.ID
	return m, nil
}
