package provider

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/jszwec/csvutil"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type CSVReader[T any] struct {
	filepath string
	// unmarshal  func([]string) (*T, error) // caller provides how to parse a row
	skipHeader bool
}

type CSVDate struct {
	time.Time
}

func NewCSVReader[T any](filepath string, skipHeader bool) *CSVReader[T] {
	return &CSVReader[T]{filepath: filepath, skipHeader: skipHeader}
}

func (r *CSVReader[T]) Read(_ context.Context) ([]*T, error) {
	var record T
	var records []*T

	f, err := os.Open(r.filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	// 1. Define the BOM-aware decoder
	// unicode.BOMOverride ensures that if a BOM is found, it's used to
	// determine the encoding and then stripped.
	win1252ToUTF8 := unicode.UTF8.NewDecoder()

	// 2. Wrap the file reader
	// This "cleanReader" will now provide a stream without the leading BOM bytes.
	cleanReader := transform.NewReader(f, unicode.BOMOverride(win1252ToUTF8))

	// 3. Use csvutil with the cleaned reader
	// We use a standard csv.Reader as the source for csvutil
	csvReader := csv.NewReader(cleanReader)

	// if the header doesn't exists
	args := []string{}
	if r.skipHeader {
		userHeader, err := csvutil.Header(record, "csv")
		if err != nil {
			return nil, err
		}
		args = userHeader
	}
	dec, err := csvutil.NewDecoder(csvReader, args...)
	if err != nil {
		return nil, err
	}
	for {
		var r T
		if err := dec.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		records = append(records, &r)
	}

	return records, nil
}

// UnmarshalCSV implements the csvutil.Unmarshaler interface
func (d *CSVDate) UnmarshalCSV(data []byte) error {
	s := string(data)
	if s == "" {
		return nil
	}

	// Define your expected CSV date format here
	// Example: "2026-03-30"
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}
