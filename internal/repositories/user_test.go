package repositories_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/online-school-mgmt/internal/models"
	repo "github.com/james-wukong/online-school-mgmt/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newAdminUserRepo(t *testing.T) (repo.AdminUserRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return repo.NewAdminUserRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleAdminUser() *models.AdminUser {
	return &models.AdminUser{
		ID:       1,
		Name:     "John",
		SchoolID: 10,
	}
}

func adminUserColumns() []string {
	// Keep in sync with all fields on models.AdminUser
	return []string{
		"id", "school_id", "name",
	}
}

func mockAdminUserRow(u *models.AdminUser) *sqlmock.Rows {
	return sqlmock.NewRows(adminUserColumns()).AddRow(
		u.ID, u.SchoolID, u.Name,
	)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestAdminUserCreate_Success(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	// Single auto-increment PK → GORM uses INSERT ... RETURNING id
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "goadmin_users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserCreate_DBError(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "goadmin_users"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserCreate_DuplicateEmail(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "goadmin_users_email_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "goadmin_users"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestAdminUserUpdate_Success(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()
	u.Name = "Updated"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "goadmin_users"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserUpdate_DBError(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "goadmin_users"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestAdminUserGetByID_Found(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockAdminUserRow(u))

	result, err := repo.GetByID(ctx, int64(u.ID))

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, u.SchoolID, result.SchoolID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetByID_NotFound(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetByID_DBError(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetByID_CancelledContext(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ── GetSchoolByUserID ─────────────────────────────────────────────────────────
// Preload("School") causes GORM to fire the School SELECT *before* the main
// query. Always match the log order: preload first, main query last.

func TestAdminUserGetSchoolByUserID_Found(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()
	sch := sampleSchool()

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockAdminUserRow(u))

	// EXPECTATION 2: School preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"`)).
		WithArgs(u.SchoolID).
		WillReturnRows(mockSchoolRow(sch))

	result, err := repo.GetSchoolByUserID(ctx, int64(u.ID))

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.School)
	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, sch.ID, result.School.ID)
	assert.Equal(t, sch.Name, result.School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetSchoolByUserID_UserNotFound(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// EXPECTATION 2: School preload
	// mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"`)).
	// 	WithArgs(int64(999)).
	// 	WillReturnRows(sqlmock.NewRows(schoolColumns()))

	result, err := repo.GetSchoolByUserID(ctx, 999)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetSchoolByUserID_SchoolPreloadError(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	dbErr := errors.New("preload failed")

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockAdminUserRow(u))

	// EXPECTATION 2: School preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"`)).
		WithArgs(u.SchoolID).
		WillReturnError(dbErr)

	result, err := repo.GetSchoolByUserID(ctx, int64(u.ID))

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetSchoolByUserID_NullSchool(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()
	u.SchoolID = 0 // user not assigned to any school

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockAdminUserRow(u))

	// EXPECTATION 2: School preload empty rows
	// mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"`)).
	// 	WithArgs(u.SchoolID).
	// 	WillReturnRows(sqlmock.NewRows(schoolColumns()))

	result, err := repo.GetSchoolByUserID(ctx, int64(u.ID))

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.School) // no school was preloaded
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdminUserGetSchoolByUserID_DBError(t *testing.T) {
	repo, mock := newAdminUserRepo(t)
	ctx := context.Background()
	u := sampleAdminUser()

	dbErr := errors.New("db timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "goadmin_users"`)).
		WithArgs(u.ID, 1).
		WillReturnError(dbErr)

	result, err := repo.GetSchoolByUserID(ctx, int64(u.ID))

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}
