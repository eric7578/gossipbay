package schedule

import (
	"encoding/json"
	"io"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/pkg/errors"
)

type Scheduler struct {
	crawler *crawler.PageCrawler
}

func NewScheduler() *Scheduler {
	s := Scheduler{
		crawler: crawler.NewPageCrawler(),
	}
	return &s
}

func Pipe(src interface{}, dest io.Writer) error {
	if bytes, err := json.Marshal(src); err != nil {
		return errors.Wrap(err, "cannot format as json")
	} else if _, err := dest.Write(bytes); err != nil {
		return errors.Wrap(err, "output failed")
	}
	return nil
}
