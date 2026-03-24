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

func newSchoolRepo(t *testing.T) (SchoolRepository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	return NewSchoolRepository(gormDB), mock
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

func sampleSchool() *models.Schools {
	return &models.Schools{
		ID:   10,
		Name: "Test School",
		Code: "TS001",
	}
}

func schoolColumns() []string {
	return []string{"id", "name", "code"}
}

func mockSchoolRow(s *models.Schools) *sqlmock.Rows {
	return sqlmock.NewRows(schoolColumns()).AddRow(s.ID, s.Name, s.Code)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestSchoolGetByID_Found(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	u := sampleSchool()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"`)).
		WithArgs(u.ID, 1).
		WillReturnRows(mockSchoolRow(u))

	result, err := repo.GetByID(ctx, u.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, u.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByID_NotFound(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// GetByID propagates the error — it does NOT swallow ErrRecordNotFound
	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByID_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByID_CancelledContext(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
