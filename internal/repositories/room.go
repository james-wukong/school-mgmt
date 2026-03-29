package repositories

import (
	"context"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, t *models.Rooms) error
	CreateWithAssoc(
		ctx context.Context, t *models.Rooms,
		tt []*models.RoomTimeslots,
	) error
	Update(ctx context.Context, t *models.Rooms) error
	UpdateRoomStatus(ctx context.Context, t *models.Rooms) error
	UpdateWithAssoc(
		ctx context.Context, t *models.Rooms,
		tt []*models.RoomTimeslots,
		semID int64,
	) error
	GetByID(ctx context.Context, id int64) (*models.Rooms, error)
}

type roomRepo struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepo{
		db: db,
	}
}

func (r *roomRepo) Create(ctx context.Context, t *models.Rooms) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *roomRepo) CreateWithAssoc(
	ctx context.Context, t *models.Rooms,
	tt []*models.RoomTimeslots,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create the Room
		if err := tx.Omit("School", "Timeslots").Create(t).Error; err != nil {
			return err
		}

		// 2. Update room-timeslot pairs (room id)
		for _, pair := range tt {
			pair.RoomID = t.ID
		}

		// 3. Insert room-timeslot pairs
		if err := tx.Model(&models.RoomTimeslots{}).Create(tt).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *roomRepo) Update(ctx context.Context, t *models.Rooms) error {
	// r.db.Model(&t).Select("IsActive").Updates(t).Error
	return r.db.WithContext(ctx).Save(t).Error
}

// UpdateRoomStatus only update is_active field for a room
func (r *roomRepo) UpdateRoomStatus(
	ctx context.Context, t *models.Rooms,
) error {
	// r.db.Model(&t).Select("IsActive").Updates(t).Error
	return r.db.WithContext(ctx).
		Model(t).
		Where("id = ?", t.ID).
		UpdateColumn("is_active", t.IsActive).
		Error
}

// UpdateWithAssoc update rooms table
// removes previously attached room-timeslot and room-timeslot pairs
// and insert new pairs
func (r *roomRepo) UpdateWithAssoc(
	ctx context.Context,
	t *models.Rooms,
	tt []*models.RoomTimeslots,
	semID int64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Upsert the Room
		if err := tx.Omit("School", "Timeslots").Save(t).Error; err != nil {
			return err
		}

		// 2 Remove the previous teacher-timeslot pair for the semester
		subQuery := tx.Model(&models.Timeslots{}).
			Select("id").
			Where("semester_id = ?", semID)
		if err := tx.Where("room_id = ?", t.ID).
			Where("timeslot_id IN (?)", subQuery).
			Delete(&models.RoomTimeslots{}).
			Error; err != nil {
			return err
		}

		// 3. Insert room-timeslot pairs
		return tx.Model(&models.RoomTimeslots{}).Create(tt).Error
	})
}

func (r *roomRepo) GetByID(ctx context.Context, id int64) (*models.Rooms, error) {
	var room models.Rooms

	err := r.db.WithContext(ctx).
		Preload("School").
		Preload("Timeslots").
		First(&room, id).Error

	if err != nil {
		return nil, err
	}

	return &room, nil
}
