package schedule

import (
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
)

type WeeklySchedule struct {
	rawWeeklySchedule

	weekday time.Weekday
}

type rawWeeklySchedule struct {
	DailySchedule

	Weekday string `yaml:"weekday"`
}

func (s WeeklySchedule) GetType() string {
	return "weekly"
}

func (s WeeklySchedule) NextUpdate() time.Time {
	nextDailyUpdate := s.DailySchedule.NextUpdate()
	thisWeeksUpdate := nextDailyUpdate.
		Add(-24 * time.Hour * time.Duration(nextDailyUpdate.Weekday())).
		Add(24 * time.Hour * time.Duration(s.weekday))

	if time.Now().Before(thisWeeksUpdate) {
		return thisWeeksUpdate
	}

	return thisWeeksUpdate.Add(7 * 24 * time.Hour)
}

var _ yaml.BytesMarshaler = (*WeeklySchedule)(nil)
var _ yaml.BytesUnmarshaler = (*WeeklySchedule)(nil)

// Custom YAML encoding
func (s *WeeklySchedule) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(&s.rawWeeklySchedule)
}

// Custom YAML decoding
func (s *WeeklySchedule) UnmarshalYAML(data []byte) error {
	captureWeekday := struct {
		Weekday string `yaml:"weekday"`
	}{}

	if err := yaml.Unmarshal(data, &captureWeekday); err != nil {
		return err
	}

	s.rawWeeklySchedule = rawWeeklySchedule{}
	if err := yaml.Unmarshal(data, &s.rawWeeklySchedule); err != nil {
		return err
	}

	s.Weekday = captureWeekday.Weekday

	parsedTime, err := time.ParseInLocation("Monday", s.Weekday, &s.tz)
	if err != nil {
		return fmt.Errorf("invalid 'weekday' %q: %w", s.Weekday, err)
	}

	s.weekday = parsedTime.Weekday()

	return nil
}
