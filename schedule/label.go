package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	scheduleLabels []string = []string{"weekly", "daily"}
)

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

func parseIssueLabels(labels []string) (time.Time, time.Time, float64) {
	var (
		from    time.Time
		to      time.Time
		deviate float64
	)
	for _, label := range labels {
		if isScheduleLabel(label) {
			from, to = getTimeRanges(label)
		} else if f, ok := isTrendingLabel(label); ok {
			deviate = f
		}
	}
	return from, to, deviate
}
