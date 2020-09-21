package flagutil

import (
	"fmt"
	"strconv"
	"time"
)

var (
	taipei *time.Location
)

func init() {
	var err error
	taipei, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
}

func ParseSchedule(schedule string) (from, to time.Time, err error) {
	now := time.Now()
	to = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)

	switch schedule {
	case "weekly":
		from = to.Add(-7 * 24 * time.Hour)
	case "daily":
		from = to.Add(-24 * time.Hour)
	}

	if from.IsZero() || to.IsZero() {
		err = fmt.Errorf("invalid schedule %s", schedule)
	}

	return from, to, err
}

func offsetDays(s string) time.Time {
	offset, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Errorf("days offset must be an int, got: %s", s))
	}

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)
	return t.Add(time.Duration(offset*24) * time.Hour)
}
