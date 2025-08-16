package schedule

import (
	"time"
)

type ASAPSchedule struct{}

func (s ASAPSchedule) GetType() string {
	return "asap"
}

func (s ASAPSchedule) NextUpdate() time.Time {
	return time.Now()
}
