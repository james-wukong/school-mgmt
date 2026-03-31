package provider

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

type TextReader[T any] struct {
	data string
}

type TextDate struct {
	time.Time
}

func NewTextReader[T any](data string) *TextReader[T] {
	return &TextReader[T]{
		data: data,
	}
}

func (r *TextReader[T]) Read(_ context.Context) ([]*T, error) {
	var results []*T
	// Pass the address (&room) so the function can modify the variable
	err := json.Unmarshal([]byte(r.data), &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// UnmarshalJSON handles the "09:00" -> HourMinute conversion
// When json.Unmarshal encounters a field of type HourMinute,
// it checks if that type has an UnmarshalJSON([]byte) error method
func (d *TextDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
