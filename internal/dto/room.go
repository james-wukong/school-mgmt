package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
)

type RoomBase struct {
	SchoolID    int64  `form:"school_id" validate:"required"`
	Code        string `form:"code" validate:"required,max=50"`
	Name        string `form:"name" validate:"required,max=100"`
	RoomType    string `form:"room_type" validate:"omitempty,max=50"`
	Capacity    int    `form:"capacity" validate:"omitempty,min=1"`
	FloorNumber int    `form:"floor_number" validate:"omitempty"`
	Building    string `form:"building" validate:"omitempty,max=100"`
	IsActive    bool   `form:"is_active"` // Pointer handles 'false' vs 'missing' (nil)

	// Many-to-Many: GoAdmin sends multi-selects as a slice of strings or IDs
	TimeslotIDs []int64 `form:"timeslots[]" validate:"omitempty,unique"`
}

type RoomCreateRequest struct {
	RoomBase
}

type RoomUpdateRequest struct {
	ID int64 `form:"id" validate:"required"` // The ID is mandatory
	RoomBase
}

type RoomStatusUpdateRequest struct {
	ID       int64 `form:"id" validate:"required"` // The ID is mandatory
	IsActive bool  `form:"is_active"`
}

func (req *RoomBase) toModel() (*models.Rooms, error) {
	return &models.Rooms{
		SchoolID:    req.SchoolID,
		Code:        req.Code,
		Name:        req.Name,
		RoomType:    models.RoomType(req.RoomType),
		Capacity:    req.Capacity,
		FloorNumber: &req.FloorNumber,
		Building:    req.Building,
		IsActive:    req.IsActive,
	}, nil
}

func (req *RoomCreateRequest) ToModel() (*models.Rooms, error) {
	return req.toModel()
}

func (req *RoomUpdateRequest) ToModel() (*models.Rooms, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}

func (req *RoomStatusUpdateRequest) ToModel() *models.Rooms {
	return &models.Rooms{
		ID:       req.ID,
		IsActive: req.IsActive,
	}
}
