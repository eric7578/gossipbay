package flagutil

import (
	"fmt"
	"strconv"
	"strings"
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

func init() {
	var err error
	taipei, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
}

func ParseDaysExpression(s string) (from, to time.Time) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)

	segs := strings.Split(s, ":")
	if segs[0] != "" {
		from = offsetDays(segs[0])
	}
	if segs[1] != "" {
		to = offsetDays(segs[1])
	}

	if to.IsZero() && startOfToday.After(from) {
		return from, startOfToday
	}
	if from.IsZero() && startOfToday.Before(to) {
		return startOfToday, to
	}
	return from, to
}

func ParseSchedule(schedule string) (time.Time, time.Time) {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)

	switch schedule {
	case "weekly":
		return to.Add(-7 * 24 * time.Hour), to
	case "daily":
		return to.Add(-24 * time.Hour), to
	default:
		panic(fmt.Errorf("invalid schedule %s", schedule))
	}
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
