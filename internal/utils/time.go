package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// UTCTime wraps time.Time to automatically handle UTC storage and local display
type UTCTime struct {
	time.Time
}

// Now returns the current time as UTCTime
func Now() UTCTime {
	return UTCTime{time.Now()}
}

// NewUTCTime creates a UTCTime from a standard time.Time
func NewUTCTime(t time.Time) UTCTime {
	return UTCTime{t}
}

// Value implements driver.Valuer for database storage
// Always stores time in UTC format
func (t UTCTime) Value() (driver.Value, error) {
	return t.UTC().Format("2006-01-02 15:04:00"), nil
}

// Scan implements sql.Scanner for database retrieval
// Reads UTC time and converts to local for display
func (t *UTCTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v.In(time.Local)
		return nil
	case string:
		parsed, err := time.ParseInLocation("2006-01-02 15:04:00", v, time.UTC)
		if err != nil {
			return err
		}
		t.Time = parsed.In(time.Local)
		return nil
	case []byte:
		parsed, err := time.ParseInLocation("2006-01-02 15:04:00", string(v), time.UTC)
		if err != nil {
			return err
		}
		t.Time = parsed.In(time.Local)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into UTCTime", value)
	}
}

// MarshalJSON for JSON serialization in local time
func (t UTCTime) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}

// UnmarshalJSON for JSON deserialization
func (t *UTCTime) UnmarshalJSON(data []byte) error {
	return t.Time.UnmarshalJSON(data)
}

// TruncateToMinute truncates time to minute precision
func (t UTCTime) TruncateToMinute() UTCTime {
	return UTCTime{t.Truncate(time.Minute)}
}

// FormatUTC formats time as UTC for database queries
func (t UTCTime) FormatUTC() string {
	return t.UTC().Format("2006-01-02 15:04:00")
}

// ToUTC converts any time.Time to UTC for database queries
func ToUTC(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:00")
}