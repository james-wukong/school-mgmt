package provider_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/james-wukong/online-school-mgmt/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── test structs ──────────────────────────────────────────────────────────────

type testPerson struct {
	Name     string           `csv:"name"`
	Age      int              `csv:"age"`
	IsActive bool             `csv:"is_active"`
	HireDate provider.CSVDate `csv:"hire_date"`
}

type testSimple struct {
	ID   int    `csv:"id"`
	Code string `csv:"code"`
}

// ── helpers ───────────────────────────────────────────────────────────────────

// writeTempCSV creates a temp file with the given content and returns its path.
// The file is automatically cleaned up when the test ends.
func writeTempCSV(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test_*.csv")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

// ── Read — with header (skipHeader: false) ────────────────────────────────────

func TestCSVReader_Read_WithHeader_Success(t *testing.T) {
	content := `name,age,is_active,hire_date
Alice,30,true,2025-01-15
Bob,25,false,2024-06-01
`
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, false)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	require.Len(t, results, 2)

	assert.Equal(t, "Alice", results[0].Name)
	assert.Equal(t, 30, results[0].Age)
	assert.True(t, results[0].IsActive)
	assert.Equal(t, 2025, results[0].HireDate.Year())
	assert.Equal(t, time.January, results[0].HireDate.Month())
	assert.Equal(t, 15, results[0].HireDate.Day())

	assert.Equal(t, "Bob", results[1].Name)
	assert.Equal(t, 25, results[1].Age)
	assert.False(t, results[1].IsActive)
	assert.Equal(t, 2024, results[1].HireDate.Year())
}

func TestCSVReader_Read_WithHeader_SingleRow(t *testing.T) {
	content := `name,age,is_active,hire_date
Alice,30,true,2025-01-15
`
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, false)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "Alice", results[0].Name)
}

func TestCSVReader_Read_WithHeader_EmptyBody(t *testing.T) {
	// Header only — no data rows
	content := "name,age,is_active,hire_date\n"
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, false)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	assert.Empty(t, results)
}

// ── Read — without header (skipHeader: true) ──────────────────────────────────
// skipHeader:true means the CSV has NO header row —
// CSVReader derives the header from the struct tags instead

func TestCSVReader_Read_SkipHeader_Success(t *testing.T) {
	// No header line — data starts immediately
	content := `Alice,30,true,2025-01-15
Bob,25,false,2024-06-01
`
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, true)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "Alice", results[0].Name)
	assert.Equal(t, "Bob", results[1].Name)
}

func TestCSVReader_Read_SkipHeader_SingleRow(t *testing.T) {
	content := "Jane,40,true,2023-09-01\n"
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, true)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "Jane", results[0].Name)
	assert.Equal(t, 40, results[0].Age)
}

// ── Read — date handling ──────────────────────────────────────────────────────

func TestCSVReader_Read_EmptyDate(t *testing.T) {
	// Empty hire_date should not error — CSVDate.UnmarshalCSV handles ""
	content := `name,age,is_active,hire_date
Alice,30,true,
`
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, false)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].HireDate.IsZero())
}

func TestCSVReader_Read_DateFormats(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		wantYear int
		wantMon  time.Month
		wantDay  int
	}{
		{"standard", "2025-01-15", 2025, time.January, 15},
		{"end of year", "2024-12-31", 2024, time.December, 31},
		{"leap day", "2024-02-29", 2024, time.February, 29},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content := "name,age,is_active,hire_date\nAlice,30,true," + tc.dateStr + "\n"
			path := writeTempCSV(t, content)
			reader := provider.NewCSVReader[testPerson](path, false)

			results, err := reader.Read(context.Background())

			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, tc.wantYear, results[0].HireDate.Year())
			assert.Equal(t, tc.wantMon, results[0].HireDate.Month())
			assert.Equal(t, tc.wantDay, results[0].HireDate.Day())
		})
	}
}

// ── Read — file errors ────────────────────────────────────────────────────────

func TestCSVReader_Read_FileNotFound(t *testing.T) {
	reader := provider.NewCSVReader[testPerson](
		filepath.Join(t.TempDir(), "nonexistent.csv"),
		false,
	)

	results, err := reader.Read(context.Background())

	assert.Nil(t, results)
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestCSVReader_Read_EmptyFile(t *testing.T) {
	path := writeTempCSV(t, "")
	reader := provider.NewCSVReader[testPerson](path, false)

	results, err := reader.Read(context.Background())

	// csvutil returns EOF on empty file — no records, no error
	require.ErrorIs(t, err, io.EOF)
	assert.Empty(t, results)
}

// ── Read — multiple rows pointer independence ─────────────────────────────────

func TestCSVReader_Read_RowsAreIndependentPointers(t *testing.T) {
	content := `id,code
1,AAA
2,BBB
3,CCC
`
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testSimple](path, false)

	results, err := reader.Read(context.Background())

	require.NoError(t, err)
	require.Len(t, results, 3)

	// Mutating one result must not affect others
	results[0].Code = "MODIFIED"
	assert.Equal(t, "BBB", results[1].Code)
	assert.Equal(t, "CCC", results[2].Code)
}

// ── Read — context (documents current behaviour) ─────────────────────────────

func TestCSVReader_Read_ContextIsUnused(t *testing.T) {
	// The current implementation ignores ctx (_ context.Context).
	// This test documents that behaviour — a cancelled context does NOT
	// stop the read. If cancellation support is added later, update this test.
	content := `name,age,is_active,hire_date
Alice,30,true,2025-01-15
`
	path := writeTempCSV(t, content)
	reader := provider.NewCSVReader[testPerson](path, false)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancelled before Read is called

	results, err := reader.Read(ctx)

	// Still succeeds because ctx is not checked
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

// ── CSVDate.UnmarshalCSV ──────────────────────────────────────────────────────

func TestCSVDate_UnmarshalCSV_ValidDate(t *testing.T) {
	var d provider.CSVDate
	err := d.UnmarshalCSV([]byte("2025-06-15"))

	require.NoError(t, err)
	assert.Equal(t, 2025, d.Year())
	assert.Equal(t, time.June, d.Month())
	assert.Equal(t, 15, d.Day())
}

func TestCSVDate_UnmarshalCSV_EmptyString(t *testing.T) {
	var d provider.CSVDate
	err := d.UnmarshalCSV([]byte(""))

	require.NoError(t, err)
	assert.True(t, d.IsZero())
}

func TestCSVDate_UnmarshalCSV_InvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"wrong separator", "2025/06/15"},
		{"DD-MM-YYYY", "15-06-2025"},
		{"text", "not-a-date"},
		{"partial", "2025-06"},
		{"US format", "06/15/2025"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var d provider.CSVDate
			err := d.UnmarshalCSV([]byte(tc.input))
			assert.Error(t, err)
		})
	}
}
