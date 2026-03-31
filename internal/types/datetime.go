package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CivilDate time.Time

const dateFormat = "2006-01-02"

// UnmarshalJSON JSON Support
func (c *CivilDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		return nil
	}
	t, err := time.Parse(`"`+dateFormat+`"`, s)
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
	t, err := time.Parse(dateFormat, string(data))
	if err != nil {
		return err
	}
	*c = CivilDate(t)
	return nil
}

// Value GORM/SQL Support (So you can save it directly to the DB)
func (c CivilDate) Value() (driver.Value, error) {
	return time.Time(c).Format(dateFormat), nil
}

func (c *CivilDate) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		*c = CivilDate(t)
		return nil
	}
	return fmt.Errorf("cannot scan %v into CivilDate", value)
}
