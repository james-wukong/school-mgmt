package repositories

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTeacherRepo(t *testing.T) (TeacherRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return NewTeacherRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleTeacher() *models.Teachers {
	return &models.Teachers{
		ID:        1,
		FirstName: "John",
		SchoolID:  10,
	}
}

func teacherColumns() []string {
	// Keep in sync with all fields on models.Teachers
	return []string{
		"id", "school_id", "first_name",
	}
}

func mockTeacherRow(u *models.Teachers) *sqlmock.Rows {
	return sqlmock.NewRows(teacherColumns()).AddRow(
		u.ID, u.SchoolID, u.FirstName,
	)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestTeacherCreate_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()

	// Single auto-increment PK → GORM uses INSERT ... RETURNING id
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherCreate_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherCreate_DuplicateEmail(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "teachers_email_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestTeacherUpdate_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()
	u.FirstName = "Updated"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherUpdate_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestTeacherGetByID_Found(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()
	sch := sampleSchool()

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockTeacherRow(u))

	// EXPECTATION 2: School preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"`)).
		WithArgs(u.SchoolID).
		WillReturnRows(mockSchoolRow(sch))

	result, err := repo.GetByID(ctx, u.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.School)
	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, u.SchoolID, result.SchoolID)
	assert.Equal(t, sch.ID, result.School.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByID_NotFound(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByID_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByID_CancelledContext(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
