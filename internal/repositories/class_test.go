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

func newClassRepo(t *testing.T) (ClassRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return NewClassRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleClass() *models.Classes {
	return &models.Classes{
		ID:           2,
		SemesterID:   3,
		Grade:        2,
		ClassName:    "A",
		StudentCount: 30,
	}
}

func classColumns() []string {
	// Keep in sync with all fields on models.Classes
	return []string{
		"id", "school_id", "semester_id", "grade", "class", "student_count",
	}
}

func mockClassRow(u *models.Classes) *sqlmock.Rows {
	return sqlmock.NewRows(classColumns()).AddRow(
		u.ID, u.SchoolID, u.SemesterID, u.Grade, u.ClassName, u.StudentCount,
	)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestClassCreate_Success(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()

	// Single auto-increment PK → GORM uses INSERT ... RETURNING id
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassCreate_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassCreate_Duplicate(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "classes_email_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassCreate_CancelledContext(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()
	// cancel()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnError(context.Canceled)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestClassUpdate_Success(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()
	u.Grade = 10

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassUpdate_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── UpdateWithSemester ────────────────────────────────────────────────────────
// Select("Semester").Save(t) → UPDATE classes + UPDATE semesters in one tx

func TestClassUpdateWithSemester_Success(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	c := sampleClass()
	c.Grade = 5
	c.Semester = sampleSemester()
	c.Semester.Semester = 4

	mock.ExpectBegin()

	// STEP 1: tx.Save(t) upserts the populated Semester association
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WithArgs(
			c.Semester.SchoolID,
			c.Semester.Year,
			c.Semester.Semester,
			c.Semester.StartDate,
			c.Semester.EndDate,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(c.Semester.ID))

	// STEP 2: tx.Save(t) then updates the classes row
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.UpdateWithSemester(ctx, c)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestClassUpdateWithSemester_SemesterUpdateError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	c := sampleClass()
	c.Semester = sampleSemester()

	dbErr := errors.New("semester update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters" SET`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.UpdateWithSemester(ctx, c)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassUpdateWithSemester_SemesterUpsertError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	c := sampleClass()
	c.Semester = sampleSemester()

	dbErr := errors.New("semester upsert failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.UpdateWithSemester(ctx, c)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassUpdateWithSemester_ClassUpdateError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	c := sampleClass()
	c.Semester = sampleSemester()

	dbErr := errors.New("class update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(c.Semester.ID))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes" SET`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.UpdateWithSemester(ctx, c)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────-

func TestClassDelete_Success(t *testing.T) {
	// Setup repo and mock
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()

	// GORM DELETE by struct typically generates:
	// DELETE FROM "classes" WHERE "classes"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "classes" WHERE "classes"."id" = $1`)).
		WithArgs(u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 0 lastInsertID, 1 row affected
	mock.ExpectCommit()

	err := repo.Delete(ctx, u)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassDelete_DBError(t *testing.T) {
	// Setup repo and mock
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()

	dbErr := errors.New("permission denied")
	// GORM DELETE by struct typically generates:
	// DELETE FROM "classes" WHERE "classes"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "classes" WHERE`)).
		WithArgs(u.ID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, u)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassDelete_NonExistentRecord(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	c := sampleClass()
	c.ID = 9999

	// DELETE on a missing row — DB returns 0 rows affected, GORM does not error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "classes"`)).
		WithArgs(c.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, c)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestClassGetByID_Found(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	u := sampleClass()
	sem := sampleSemester()

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockClassRow(u))

	// EXPECTATION 2: School preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE `)).
		WithArgs(u.SemesterID).
		WillReturnRows(mockSemesterRow(sem))

	result, err := repo.GetByID(ctx, u.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Semester)
	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, u.SemesterID, result.SemesterID)
	assert.Equal(t, sem.ID, result.Semester.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetByID_NotFound(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetByID_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetByID_CancelledContext(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
