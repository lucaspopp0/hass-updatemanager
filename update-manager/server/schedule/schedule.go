package schedule

import (
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
)

type Schedule interface {
	GetType() string
	NextUpdate() time.Time
}

func UnmarshalYAML(data []byte) (Schedule, error) {
	typed := struct {
		Type string `json:"type"`
	}{}

	if err := yaml.Unmarshal(data, &typed); err != nil {
		return nil, err
	}

	switch typed.Type {
	case "asap":
		return ASAPSchedule{}, nil
	case "daily":
		daily := DailySchedule{}
		if err := daily.UnmarshalYAML(data); err != nil {
			return nil, err
		}

		return daily, nil
	case "weekly":
		weekly := WeeklySchedule{}
		if err := weekly.UnmarshalYAML(data); err != nil {
			return nil, err
		}

		return weekly, nil
	default:
		return nil, fmt.Errorf("unknown type %q", typed.Type)
	}
}
