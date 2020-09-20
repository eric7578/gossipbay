package repo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/eric7578/gossipbay/flagutil"
)

func parseTrendingOption(issue GithubIssue) crawler.TrendingOption {
	opt := crawler.TrendingOption{
		Board: issue.Title,
	}
	for _, label := range issue.Labels {
		if f, ok := isDeviateLabel(label.Name); ok {
			opt.Deviate = f
		} else if d, ok := isTimeoutLabel(label.Name); ok {
			opt.Timeout = d
		} else if from, to, err := flagutil.ParseSchedule(label.Name); err == nil {
			opt.From = from
			opt.To = to
		}
	}
	return opt
}

func isTimeoutLabel(l string) (time.Duration, bool) {
	if strings.Index(l, "timeout-") != 0 {
		return time.Duration(0), false
	}

	segs := strings.Split(l, "-")
	i, err := strconv.ParseInt(segs[1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid value for timeout tag: %s", l))
	}
	return time.Duration(i) * time.Second, true
}

func isDeviateLabel(l string) (float64, bool) {
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
