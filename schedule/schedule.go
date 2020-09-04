package schedule

import (
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

type RunOption struct {
	Board   string
	From    time.Time
	To      time.Time
	Timeout time.Duration
	Deviate float64
}

func (opt RunOption) isValid() bool {
	return opt.Board != "" && !opt.From.IsZero()
}

type BoardReport struct {
	RunOption
	Total   int
	Threads []crawler.Thread
}

type Scheduler struct {
	crawler *crawler.PageCrawler
}

func NewScheduler() *Scheduler {
	s := Scheduler{
		crawler: crawler.NewPageCrawler(),
	}
	return &s
}
