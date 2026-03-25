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

func newSemesterRepo(t *testing.T) (SemesterRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return NewSemesterRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleSemester() *models.Semesters {
	start, _ := time.Parse(models.TimeDateLayout, "2025-01-15")
	return &models.Semesters{
		ID:        3,
		SchoolID:  10,
		Year:      2025,
		Semester:  1,
		StartDate: start,
	}
}

func semesterColumns() []string {
	// Keep in sync with all fields on models.Semesters
	return []string{
		"id", "school_id", "year", "semester", "start_date",
	}
}

func mockSemesterRow(u *models.Semesters) *sqlmock.Rows {
	return sqlmock.NewRows(semesterColumns()).AddRow(
		u.ID, u.SchoolID, u.Year, u.Semester, u.StartDate,
	)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestSemesterCreate_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()

	u.Classes = []*models.Classes{
		{SemesterID: u.ID, Grade: 5, ClassName: "A", StudentCount: 30},
		{SemesterID: u.ID, Grade: 5, ClassName: "B", StudentCount: 28},
	}

	// Single auto-increment PK → GORM uses INSERT ... RETURNING id
	mock.ExpectBegin()
	// 1. Expect Parent (Semester) Insert
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))

	// 2. Expect Children (Classes) Bulk Upsert
	// Note: GORM uses Query because of RETURNING "id"
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WithArgs(
			u.ID, 5, "A", 30, // First class
			u.ID, 5, "B", 28, // Second class
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2).AddRow(3))
	mock.ExpectCommit()

	err := repo.Create(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterCreate_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterCreate_DuplicateEmail(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "semesters_email_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, u)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestSemesterUpdateWithClasses_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()
	u.Classes = []*models.Classes{
		{ID: 2, SemesterID: u.ID, Grade: 5, ClassName: "A", StudentCount: 30},
		{ID: 3, SemesterID: u.ID, Grade: 5, ClassName: "B", StudentCount: 28},
	}
	u.Year = 2000

	mock.ExpectBegin()
	// STEP 1: Parent Update (Matches your latest log)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters" SET`)).
		WithArgs(
			u.SchoolID,
			u.Year, // 2000
			u.Semester,
			u.StartDate,
			u.EndDate,
			u.ID, // WHERE id = 3
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// STEP 2: Children Bulk Upsert
	// Note: Use ExpectQuery because of the RETURNING "id" clause
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WithArgs(
			u.ID, u.Classes[0].Grade, u.Classes[0].ClassName, u.Classes[0].StudentCount,
			u.ID, u.Classes[1].Grade, u.Classes[1].ClassName, u.Classes[1].StudentCount,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2).AddRow(3))

	mock.ExpectCommit()

	err := repo.UpdateWithClasses(ctx, u)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterUpdateWithClasses_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.UpdateWithClasses(ctx, u)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── UpdateWithClassAssocReplace ───────────────────────────────────────────────────-

func TestSemesterUpdateWithClassAssocReplace_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	sem := sampleSemester()
	sem.Classes = []*models.Classes{
		{ID: 2, SemesterID: sem.ID, Grade: 5, ClassName: "A", StudentCount: 30},
		{ID: 3, SemesterID: sem.ID, Grade: 5, ClassName: "B", StudentCount: 28},
	}

	mock.ExpectBegin()

	// 1. Parent Update
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters" SET`)).
		WithArgs(sem.SchoolID, sem.Year, sem.Semester, sem.StartDate, sem.EndDate, sem.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 2. Bulk Upsert of the current slice
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WithArgs(
			sem.ID, sem.Classes[0].Grade, sem.Classes[0].ClassName, sem.Classes[0].StudentCount,
			sem.ID, sem.Classes[1].Grade, sem.Classes[1].ClassName, sem.Classes[1].StudentCount,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2).AddRow(3))

	// 3. The Unlink (The fix is the arguments list here)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes" SET "semester_id"=$1 WHERE "classes"."id" NOT IN ($2,$3) AND "classes"."semester_id" = $4`)).
		WithArgs(
			nil,    // $1: The NULL value
			2,      // $2: First excluded ID
			3,      // $3: Second excluded ID
			sem.ID, // $4: The parent ID
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	err := repo.UpdateWithClassAssocReplace(ctx, sem)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── AppendClasses ───────────────────────────────────────────────────────────────────-

func TestSemesterAppendClasses_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()
	var cls []*models.Classes
	cls = append(cls, &models.Classes{
		SemesterID:   3,
		Grade:        2,
		ClassName:    "A",
		StudentCount: 30,
	})

	// GORM Association Append usually starts its own transaction internally
	mock.ExpectBegin()
	// Expect the Insert for the new class with the SemesterID pre-filled by GORM
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WithArgs(cls[0].SemesterID, cls[0].Grade, cls[0].ClassName, cls[0].StudentCount).
		WillReturnRows(sqlmock.NewRows([]string{"grade"}).AddRow(2))
	mock.ExpectCommit()

	err := repo.AppendClasses(ctx, u, cls)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterAppendClasses_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()
	var cls []*models.Classes
	cls = append(cls, &models.Classes{
		SemesterID:   3,
		Grade:        2,
		ClassName:    "A",
		StudentCount: 30,
	})

	dbErr := errors.New("insert failed")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.AppendClasses(ctx, u, cls)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────-

func TestSemesterDelete_Success(t *testing.T) {
	// Setup repo and mock
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()

	// GORM DELETE by struct typically generates:
	// DELETE FROM "semesters" WHERE "semesters"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "semesters" WHERE "semesters"."id" = $1`)).
		WithArgs(u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 0 lastInsertID, 1 row affected
	mock.ExpectCommit()

	err := repo.Delete(ctx, u)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterDelete_DBError(t *testing.T) {
	// Setup repo and mock
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()

	dbErr := errors.New("permission denied")
	// GORM DELETE by struct typically generates:
	// DELETE FROM "semesters" WHERE "semesters"."id" = $1
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "semesters" WHERE`)).
		WithArgs(u.ID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, u)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestSemesterGetByID_Found(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	u := sampleSemester()
	sch := sampleSchool()
	cls := sampleClass()

	// EXPECTATION 1: Main goadmin_users query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockSemesterRow(u))
	// EXPECTATION 2: Classes preload
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE "classes"`)).
		WithArgs(u.ID).
		WillReturnRows(mockClassRow(cls))

	// EXPECTATION 3: School preload
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

func TestSemesterGetByID_NotFound(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterGetByID_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterGetByID_CancelledContext(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// ── List ───────────────────────────────────────────────────────────────────

func TestSemesterList_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	schoolID := int64(10)
	limit := 10

	// Define the expected rows to be returned by the mock
	rows := sqlmock.NewRows(semesterColumns()).
		AddRow(3, schoolID, 2025, 1, time.Now()).
		AddRow(2, schoolID, 2024, 2, time.Now())

	// 1. ExpectQuery because we are fetching data (SELECT)
	// Note: Order matters! WHERE -> ORDER BY -> LIMIT
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "semesters" WHERE school_id = $1 ORDER BY id DESC LIMIT $2`,
	)).
		WithArgs(schoolID, limit).
		WillReturnRows(rows)

	// Execute the function
	results, err := repo.List(ctx, schoolID, limit)

	// Assertions
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(3), results[0].ID)
	assert.Equal(t, schoolID, results[0].SchoolID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_Empty(t *testing.T) {
	repo, mock := newSemesterRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(10, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	results, err := repo.List(context.Background(), 10, 10)

	assert.NoError(t, err)
	assert.Empty(t, results)
}
