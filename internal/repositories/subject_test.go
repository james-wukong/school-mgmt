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

func newSubjectRepo(t *testing.T) (SubjectRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return NewSubjectRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleSubject() *models.Subjects {
	return &models.Subjects{
		ID:          2,
		SchoolID:    10,
		Name:        "Math",
		Code:        "Math",
		RequiresLab: false,
		IsHeavy:     false,
	}
}

func subjectColumns() []string {
	// Keep in sync with all fields on models.Subjects
	return []string{
		"id", "school_id", "name", "code", "requires_lab", "is_heavy",
	}
}

func mockSubjectRow(u *models.Subjects) *sqlmock.Rows {
	return sqlmock.NewRows(subjectColumns()).AddRow(
		u.ID, u.SchoolID, u.Name, u.Code, u.RequiresLab, u.IsHeavy,
	)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestSubjectCreate_Success(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()

	// Single auto-increment PK → GORM uses INSERT ... RETURNING id
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectCreate_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectCreate_Duplicate(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "subjects_email_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectCreate_CancelledContext(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()
	// cancel()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnError(context.Canceled)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestSubjectUpdate_Success(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()
	u.Name = "Maths"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subjects"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectUpdate_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subjects"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────-

func TestSubjectDelete_Success(t *testing.T) {
	// Setup repo and mock
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()

	// GORM DELETE by struct typically generates:
	// DELETE FROM "subjects" WHERE "subjects"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subjects" WHERE "subjects"."id" = $1`)).
		WithArgs(u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 0 lastInsertID, 1 row affected
	mock.ExpectCommit()

	err := repo.Delete(ctx, u)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectDelete_DBError(t *testing.T) {
	// Setup repo and mock
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()

	dbErr := errors.New("permission denied")
	// GORM DELETE by struct typically generates:
	// DELETE FROM "subjects" WHERE "subjects"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subjects" WHERE`)).
		WithArgs(u.ID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, u)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectDelete_NonExistentRecord(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	c := sampleSubject()
	c.ID = 9999

	// DELETE on a missing row — DB returns 0 rows affected, GORM does not error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subjects"`)).
		WithArgs(c.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, c)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestSubjectGetByID_Found(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	u := sampleSubject()
	sch := sampleSchool()
	tch := sampleTeacher()

	// EXPECTATION 1: Main table query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockSubjectRow(u))

	// EXPECTATION 2: School preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE `)).
		WithArgs(u.SchoolID).
		WillReturnRows(mockSchoolRow(sch))

	// EXPECTATION 3: Composite table
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE `)).
		WithArgs(u.ID).
		WillReturnRows(sqlmock.NewRows([]string{"teacher_id", "subject_id"}).AddRow(1, 2))

	// EXPECTATION 4: Teacher preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE `)).
		WithArgs(tch.ID).
		WillReturnRows(mockTeacherRow(tch))

	result, err := repo.GetByID(ctx, u.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.School)
	require.NotNil(t, result.Teachers)

	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, u.SchoolID, result.School.ID)
	assert.Equal(t, tch.ID, result.Teachers[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByID_NotFound(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByID_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByID_CancelledContext(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
