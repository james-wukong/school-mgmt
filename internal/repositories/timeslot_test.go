package repositories_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	repo "github.com/james-wukong/online-school-mgmt/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTimeslotRepo(t *testing.T) (repo.TimeslotRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return repo.NewTimeslotRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleTimeslot() *models.Timeslots {
	start, _ := time.Parse(models.TimeSlotLayout, "09:00")
	return &models.Timeslots{
		ID:         1,
		SchoolID:   10,
		SemesterID: 3,
		DayOfWeek:  1,
		StartTime:  start,
	}
}

func timeslotColumns() []string {
	// Keep in sync with all fields on models.Timeslots
	return []string{
		"id", "school_id", "semester_id", "day_of_week", "start_time",
	}
}

func mockTimeslotRow(u *models.Timeslots) *sqlmock.Rows {
	return sqlmock.NewRows(timeslotColumns()).AddRow(
		u.ID, u.SchoolID, u.SemesterID, u.DayOfWeek, u.StartTime,
	)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestTimeslotCreate_Success(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	// Single auto-increment PK → GORM uses INSERT ... RETURNING id
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "timeslots"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotCreate_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "timeslots"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotCreate_DuplicateEmail(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "timeslots_email_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "timeslots"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestTimeslotUpdate_Success(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()
	u.DayOfWeek = 7

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "timeslots"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotUpdate_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "timeslots"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestTimeslotGetByID_Found(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockTimeslotRow(u))

	result, err := repo.GetByID(ctx, u.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, u.SchoolID, result.SchoolID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetByID_NotFound(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetByID_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetByID_CancelledContext(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ── Delete ───────────────────────────────────────────────────────────────────-

func TestTimeslotDelete_Success(t *testing.T) {
	// Setup repo and mock
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	// GORM DELETE by struct typically generates:
	// DELETE FROM "timeslots" WHERE "timeslots"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "timeslots" WHERE "timeslots"."id" = $1`)).
		WithArgs(u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 0 lastInsertID, 1 row affected
	mock.ExpectCommit()

	err := repo.Delete(ctx, u)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotDelete_DBError(t *testing.T) {
	// Setup repo and mock
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	u := sampleTimeslot()

	dbErr := errors.New("permission denied")
	// GORM DELETE by struct typically generates:
	// DELETE FROM "timeslots" WHERE "timeslots"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "timeslots" WHERE`)).
		WithArgs(u.ID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, u)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotDelete_NonExistentRecord(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	c := sampleTimeslot()
	c.ID = 9999

	// DELETE on a missing row — DB returns 0 rows affected, GORM does not error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "timeslots"`)).
		WithArgs(c.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, c)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
