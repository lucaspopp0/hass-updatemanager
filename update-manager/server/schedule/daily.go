package schedule

import (
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
)

type DailySchedule struct {
	rawDailySchedule

	time time.Duration `json:"-"`
	tz   time.Location `json:"-"`
}

type rawDailySchedule struct {
	Time string `json:"time"`
	TZ   string `json:"tz"`
}

func (s DailySchedule) GetType() string {
	return "daily"
}

func (s DailySchedule) NextUpdate() time.Time {
	now := time.Now().In(&s.tz).Truncate(time.Minute)

	todayStart := now.
		Add(-time.Hour * time.Duration(now.Hour())).
		Add(-time.Minute * time.Duration(now.Minute()))

	todaysUpdate := todayStart.Add(s.time)

	if now.Before(todaysUpdate) {
		return todaysUpdate
	}

	return todaysUpdate.Add(24 * time.Hour)
}

var _ yaml.BytesMarshaler = (*DailySchedule)(nil)
var _ yaml.BytesUnmarshaler = (*DailySchedule)(nil)

// Custom YAML encoding
func (s *DailySchedule) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(s.rawDailySchedule)
}

// Custom YAML decoding
func (s *DailySchedule) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal(data, &s.rawDailySchedule)
	if err != nil {
		return err
	}

	tz, err := time.LoadLocation(s.TZ)
	if err != nil {
		return fmt.Errorf("invalid 'tz' %q: %w", s.TZ, err)
	}

	s.tz = *tz

	parsedTime, err := time.ParseInLocation("3:04pm", s.Time, tz)
	if err != nil {
		return fmt.Errorf("invalid 'time' %q: %w", s.Time, err)
	}

	s.time = (time.Duration(parsedTime.Hour()) * time.Hour) +
		(time.Duration(parsedTime.Minute()) * time.Minute)

	return nil
}
