package schedule

import (
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

type RunOption struct {
	Board   string        `json:"board"`
	From    time.Time     `json:"from"`
	To      time.Time     `json:"to"`
	Timeout time.Duration `json:"-"`
	Deviate float64       `json:"-"`
}

func (opt RunOption) isValid() bool {
	return opt.Board != "" && !opt.From.IsZero()
}

type BoardReport struct {
	RunOption
	Total   int              `json:"total"`
	Threads []crawler.Thread `json:"threads"`
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
