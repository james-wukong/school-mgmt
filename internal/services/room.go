package services

import (
	"context"

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
) error {
	return s.repo.UpdateWithAssoc(ctx, t, rt)
}

func (s *RoomService) GetRoom(ctx context.Context, id int64) (*models.Rooms, error) {
	return s.repo.GetByID(ctx, id)
}
