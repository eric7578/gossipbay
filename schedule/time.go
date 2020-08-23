package schedule

import (
	"fmt"
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

func getTimeRanges(trending string) (time.Time, time.Time) {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)

	switch trending {
	case "weekly":
		return to.Add(-7 * 24 * time.Hour), to
	case "daily":
		return to.Add(-24 * time.Hour), to
	default:
		panic(fmt.Errorf("invalid period %s", trending))
	}
}

func offsetDays(d int64) time.Time {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, taipei)
	return t.Add(time.Duration(d*24) * time.Hour)
}
