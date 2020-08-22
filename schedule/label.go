package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	taipei         *time.Location
	scheduleLabels []string = []string{"weekly", "daily"}
)

func init() {
	var err error
	taipei, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
}

func isScheduleLabel(l string) bool {
	for _, label := range scheduleLabels {
		if label == l {
			return true
		}
	}
	return false
}

func isTrendingLabel(l string) (float64, bool) {
	if strings.Index(l, "trending-") != 0 {
		return 0.0, false
	}

	segs := strings.Split(l, "-")
	f, err := strconv.ParseFloat(segs[1], 64)
	if err != nil {
		panic(fmt.Errorf("invalid value for trending tag: %s", l))
	}

	return f, true
}

type scheduleOption struct {
	deviate float64
	from    time.Time
	to      time.Time
}

func newScheduleOption(labels []string) scheduleOption {
	opt := scheduleOption{}
	for _, label := range labels {
		if isScheduleLabel(label) {
			opt.from, opt.to = getTimeRanges(label)
		} else if f, ok := isTrendingLabel(label); ok {
			opt.deviate = f
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
