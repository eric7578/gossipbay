package crawler

import (
	"errors"

	"github.com/eric7578/gossipbay/crawler/pttweb"
)

type Worker interface {
	Run(args map[string]string) (interface{}, error)
}

type Crawler struct {
	workers map[string]Worker
	sem     chan struct{}
}

func NewCrawler() *Crawler {
	return &Crawler{
		workers: map[string]Worker{
			"pttweb": pttweb.NewPttWorker(),
		},
		sem: make(chan struct{}, 10),
	}
}

func (c *Crawler) CreateJob(t string, args map[string]string) (interface{}, error) {
	c.sem <- struct{}{}
	defer func() {
		<-c.sem
	}()

	if w, ok := c.workers[t]; !ok {
		return nil, errors.New("invalid crawler type")
	} else {
		return w.Run(args)
	}
}
