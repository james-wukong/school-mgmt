package repositories

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

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
		IsActive:  true,
	}
}

func teacherColumns() []string {
	// Keep in sync with all fields on models.Teachers
	return []string{
		"id", "school_id", "first_name", "is_active",
	}
}

func mockTeacherRow(u *models.Teachers) *sqlmock.Rows {
	return sqlmock.NewRows(teacherColumns()).AddRow(
		u.ID, u.SchoolID, u.FirstName, u.IsActive,
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

// ── UpdateTeacherStatus ────────────────────────────────────────────────────────────────────

func TestTeacherUpdateTeacherStatus_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()
	u.IsActive = false

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateTeacherStatus(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherUpdateTeacherStatus_DBError(t *testing.T) {
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

// ── UpdateWithTeacherSubject ────────────────────────────────────────────────────────────────────

func TestTeacherUpdateWithTeacherSubject_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	u := sampleTeacher()
	u.IsActive = false
	u.Subjects = []*models.Subjects{
		{ID: 1, SchoolID: 10, Name: "Mathematics", Code: "MATH101"},
		{ID: 2, SchoolID: 10, Name: "Physics", Code: "PHYS101"},
	}
	ts := []*models.TeacherSubjects{
		{TeacherID: 1, SubjectID: 1},
		{TeacherID: 1, SubjectID: 2},
	}

	mock.ExpectBegin()
	// Expectation 1. update main: teachers table
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expectation 2. delete from teacher_subjecst where teacher id =
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "teacher_subjects" WHERE`)).
		WithArgs(u.ID).
		WillReturnResult(sqlmock.NewResult(2, 2))

	// Expectation 3. insert new teacher-subject pairs
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_subjects" ("teacher_id","subject_id")`)).
		WithArgs(
			u.ID, u.Subjects[0].ID,
			u.ID, u.Subjects[1].ID,
		).
		WillReturnResult(sqlmock.NewResult(2, 2))
	mock.ExpectCommit()

	err := repo.UpdateWithTeacherSubject(ctx, u, ts)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── ReplaceWithSubjectAssoc ───────────────────────────────────────────────────
//
// Transaction sequence for many2many Replace with Unscoped:
//   1. UPDATE "teachers"                          ← tx.Omit(...).Save(t)
//   2. DELETE FROM "teacher_subjects"             ← Unscoped Replace clears join table
//      WHERE "teacher_subjects"."teacher_id" = ?
//   3. INSERT INTO "teacher_subjects"             ← Replace inserts new join rows
//      ("teacher_id","subject_id") VALUES (?,?)

func sampleTeacher111() *models.Teachers {
	email := "jane@school.com"
	return &models.Teachers{
		ID:               1000,
		SchoolID:         10,
		EmployeeID:       1001,
		FirstName:        "Jane",
		LastName:         "Doe",
		Email:            &email,
		HireDate:         time.Now(),
		EmploymentType:   "Full-time",
		MaxClassesPerDay: 5,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func sampleSubjects() []*models.Subjects {
	return []*models.Subjects{
		{ID: 1, SchoolID: 10, Name: "Mathematics", Code: "MATH101"},
		{ID: 2, SchoolID: 10, Name: "Physics", Code: "PHYS101"},
	}
}

func TestTeacherReplaceWithSubjectAssoc_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	tch := sampleTeacher111()
	tch.Subjects = sampleSubjects()

	mock.ExpectBegin()

	// STEP 1: tx.Omit("School","Subjects").Save(t)
	// ID has <-:false so GORM uses UPDATE not INSERT
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// STEP 2: Unscoped Replace hard-deletes existing join rows
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "teacher_subjects" WHERE "teacher_subjects"."teacher_id" = $1`)).
		WithArgs(tch.ID).
		WillReturnResult(sqlmock.NewResult(2, 2))

	// STEP 3: Insert new join rows for each subject
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_subjects" ("teacher_id","subject_id")`)).
		WithArgs(
			tch.ID, tch.Subjects[0].ID,
			tch.ID, tch.Subjects[1].ID,
		).
		WillReturnResult(sqlmock.NewResult(2, 2))

	mock.ExpectCommit()

	err := repo.ReplaceWithSubjectAssoc(ctx, tch)

	assert.NoError(t, err)
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
