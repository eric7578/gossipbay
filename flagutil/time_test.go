package flagutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDaysExpression(t *testing.T) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)
	startOf7DaysAgo := startOfToday.Add(-7 * 24 * time.Hour)
	startOf7DaysLater := startOfToday.Add(7 * 24 * time.Hour)

	t1, t2 := ParseDaysExpression("-7:")
	assert.True(t, t1.Equal(startOf7DaysAgo))
	assert.True(t, t2.Equal(startOfToday))

	t1, t2 = ParseDaysExpression(":-7")
	assert.True(t, t1.IsZero())
	assert.True(t, t2.Equal(startOf7DaysAgo))

	t1, t2 = ParseDaysExpression(":+7")
	assert.True(t, t1.Equal(startOfToday))
	assert.True(t, t2.Equal(startOf7DaysLater))

	t1, t2 = ParseDaysExpression("+7:")
	assert.True(t, t1.Equal(startOf7DaysLater))
	assert.True(t, t2.IsZero())

	t1, t2 = ParseDaysExpression("-7:+7")
	assert.True(t, t1.Equal(startOf7DaysAgo))
	assert.True(t, t2.Equal(startOf7DaysLater))
}
