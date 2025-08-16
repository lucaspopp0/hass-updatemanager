package schedule_test

import (
	_ "embed"
	"testing"

	"github.com/lucaspopp0/hass-update-manager/update-manager/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testfiles/invalid-type.yml
	invalidTypeYML []byte

	//go:embed testfiles/asap.yml
	asapYML []byte

	//go:embed testfiles/daily.yml
	dailyYML []byte

	//go:embed testfiles/daily-invalid-tz.yml
	dailyInvalidTZYML []byte

	//go:embed testfiles/daily-invalid-time.yml
	dailyInvalidTimeYML []byte

	//go:embed testfiles/weekly.yml
	weeklyYML []byte
)

func TestInvalidType(t *testing.T) {
	_, err := schedule.UnmarshalYAML(invalidTypeYML)
	require.ErrorContains(t, err, "unknown type")
}

func TestASAP(t *testing.T) {
	out, err := schedule.UnmarshalYAML(asapYML)
	require.NoError(t, err)
	require.IsType(t, out, schedule.ASAPSchedule{})
}

func TestDaily(t *testing.T) {
	t.Run("valid-daily", func(t *testing.T) {
		out, err := schedule.UnmarshalYAML(dailyYML)
		require.NoError(t, err)
		require.IsType(t, out, schedule.DailySchedule{})

		nextUpdate := out.NextUpdate()
		assert.Equal(t, 9, nextUpdate.Hour())
		assert.Equal(t, 30, nextUpdate.Minute())
	})

	t.Run("invalid-tz", func(t *testing.T) {
		_, err := schedule.UnmarshalYAML(dailyInvalidTZYML)
		require.ErrorContains(t, err, "'tz'")
	})

	t.Run("invalid-time", func(t *testing.T) {
		_, err := schedule.UnmarshalYAML(dailyInvalidTimeYML)
		require.ErrorContains(t, err, "'time'")
	})
}

func TestWeekly(t *testing.T) {
	out, err := schedule.UnmarshalYAML(weeklyYML)
	require.NoError(t, err)
	require.IsType(t, out, schedule.WeeklySchedule{})

	nextUpdate := out.NextUpdate()
	assert.Equal(t, 1, nextUpdate.Weekday())
	assert.Equal(t, 9, nextUpdate.Hour())
	assert.Equal(t, 0, nextUpdate.Minute())
}
