package pttweb

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/eric7578/gossipbay/crawler/utils"
)

func parseArgs(args map[string]string) (trendingArgs, error) {
	targs := trendingArgs{
		board: args["board"],
	}
	if f, ok := parseDeviateArg(args["deviate"]); ok {
		targs.deviate = f
	}
	if d, ok := parseTimeoutArg(args["timeout"]); ok {
		targs.timeout = d
	}
	if from, to, ok := parseRangeArg(args["range"]); ok {
		targs.from = from
		targs.to = to
	}
	if targs.deviate <= float64(0) {
		return targs, errors.New("invalid deviate")
	}
	if targs.from.IsZero() || targs.to.IsZero() {
		return targs, errors.New("from/to is required")
	}
	return targs, nil
}

func parseTimeoutArg(arg string) (time.Duration, bool) {
	if arg == "" {
		return 0, true
	}

	i, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid value for timeout tag: %s", arg))
	}
	return time.Duration(i) * time.Second, true
}

func parseDeviateArg(arg string) (float64, bool) {
	f, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		panic(fmt.Errorf("invalid value for trending tag: %s", arg))
	}
	return f, true
}

func parseRangeArg(arg string) (time.Time, time.Time, bool) {
	var (
		now  = time.Now()
		from time.Time
		to   time.Time = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, utils.Taipei)
	)

	switch arg {
	case "last-week":
		return to.Add(-7 * 24 * time.Hour), to, true
	case "yesterday":
		return to.Add(-24 * time.Hour), to, true
	default:
		return from, to, false
	}
}
