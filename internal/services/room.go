package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/repositories"
	"gorm.io/gorm"
)

type RoomService struct {
	repo repositories.RoomRepository
}

func NewRoomService(db *gorm.DB) *RoomService {
	return &RoomService{
		repo: repositories.NewRoomRepository(db),
	}
}

func (s *RoomService) CreateRoom(ctx context.Context, t *models.Rooms) error {
	return s.repo.Create(ctx, t)
}

func (s *RoomService) CreateRoomsInBatches(ctx context.Context,
	buildingName string,
	schoolID int64,
	totalFloor, numOfRooms int,
) error {
	var rooms []*models.Rooms
	if totalFloor == 0 || numOfRooms == 0 {
		return errors.New("empty rooms")
	}
	for f := range totalFloor {
		for n := range numOfRooms {
			var room models.Rooms
			floor := f + 1
			room.SchoolID = schoolID
			room.Building = buildingName
			room.Capacity = 40
			room.Name = fmt.Sprintf("Room %d-%02d", floor, n+1)
			room.Code = fmt.Sprintf("R%d%02d", floor, n+1)
			room.FloorNumber = &floor
			room.RoomType = models.Regular
			room.IsActive = true
			rooms = append(rooms, &room)
		}
	}

	return s.repo.CreateInBatches(ctx, rooms)
}

func (s *RoomService) CreateWithAssoc(
	ctx context.Context, t *models.Rooms,
	rt []*models.RoomTimeslots,
) error {
	return s.repo.CreateWithAssoc(ctx, t, rt)
}

func (s *RoomService) UpdateStatus(ctx context.Context, t *models.Rooms) error {
	return s.repo.UpdateRoomStatus(ctx, t)
}

func (s *RoomService) UpdateWithAssoc(
	ctx context.Context,
	t *models.Rooms,
	rt []*models.RoomTimeslots,
	semID int64,
) error {
	return s.repo.UpdateWithAssoc(ctx, t, rt, semID)
}

func (s *RoomService) GetRoom(ctx context.Context, id int64) (*models.Rooms, error) {
	return s.repo.GetByID(ctx, id)
}
