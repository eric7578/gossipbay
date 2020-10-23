package crawler

import (
	"errors"

	"github.com/eric7578/gossipbay/crawler/pttweb"
)

type Worker interface {
	Accept(args map[string]string) bool
	Run(args map[string]string) (interface{}, error)
}

type Crawler struct {
	workers []Worker
	sem     chan struct{}
}

func NewCrawler() *Crawler {
	return &Crawler{
		workers: []Worker{
			pttweb.NewPttWorker(),
		},
		sem: make(chan struct{}, 10),
	}
}

func (c *Crawler) CreateJob(args map[string]string) (interface{}, error) {
	c.sem <- struct{}{}
	defer func() {
		<-c.sem
	}()

	var w Worker
	for _, wr := range c.workers {
		if wr.Accept(args) {
			w = wr
			break
		}
	}

	if w == nil {
		return nil, errors.New("invalid crawler type")
	}

	return w.Run(args)
}
