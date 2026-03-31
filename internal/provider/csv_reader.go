package provider

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/jszwec/csvutil"
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
	csvReader := csv.NewReader(f)

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
