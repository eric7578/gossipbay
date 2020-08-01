package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var taipei *time.Location

func init() {
	var err error
	taipei, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
}

type scheduleOption struct {
	deviate float64
	from    time.Time
	to      time.Time
}

func newScheduleOption(labels []string) scheduleOption {
	opt := scheduleOption{}

	for _, label := range labels {
		segs := strings.Split(label, "-")
		switch segs[0] {
		case "deviate":
			f, err := strconv.ParseFloat(segs[1], 64)
			if err != nil {
				panic(fmt.Errorf("invalid value for deviate tag: %s", label))
			}
			opt.deviate = f
		case "trending":
			opt.from, opt.to = getTimeRanges(segs[1])
		}
	}
	return opt
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
