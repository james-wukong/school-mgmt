package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CivilDate time.Time

const dateFormat = "2006-01-02"

// UnmarshalJSON JSON Support
func (c *CivilDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	*c = CivilDate(t)
	return nil
}

// UnmarshalCSV CSV Support (csvutil uses UnmarshalCSV or TextUnmarshaler)
func (c *CivilDate) UnmarshalCSV(data []byte) error {
	return c.UnmarshalText(data)
}

// UnmarshalText Form/Text Support (Used by many Go form decoders)
func (c *CivilDate) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	t, err := time.Parse(dateFormat, string(data))
	if err != nil {
		return err
	}
	*c = CivilDate(t)
	return nil
}

// Value GORM/SQL Support (So you can save it directly to the DB)
func (c CivilDate) Value() (driver.Value, error) {
	if time.Time(c).IsZero() {
		return nil, nil
	}
	return time.Time(c).Format(dateFormat), nil
}

func (c *CivilDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*c = CivilDate(v)
	case string:
		t, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
		*c = CivilDate(t)
	case []byte:
		t, err := time.Parse(dateFormat, string(v))
		if err != nil {
			return err
		}
		*c = CivilDate(t)
	default:
		return fmt.Errorf("cannot scan %T into CivilDate", value)
	}

	return nil
}

type Int64Slice []int64

// UnmarshalCSV converts a comma-separated string "1,2,3" into []int64
func (s *Int64Slice) UnmarshalCSV(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Split by comma
	parts := strings.Split(string(data), ",")
	for _, p := range parts {
		val, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			return err
		}
		*s = append(*s, val)
	}
	return nil
}
