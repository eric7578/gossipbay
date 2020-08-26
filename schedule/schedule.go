package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/eric7578/gossipbay/repo"
)

type issueConfig struct {
	deviate float64
}

func issueConfigFromLabels(labels []string) issueConfig {
	cfg := issueConfig{}
	for _, label := range labels {
		if f, ok := isTrendingLabel(label); ok {
			cfg.deviate = f
		}
	}
	return cfg
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

func Run(r repo.Repository, label string, from, to time.Time) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("schedule failed %s", e)
		}
	}()

	var wg sync.WaitGroup
	issues := r.ListIssues(label)
	wg.Add(len(issues))
	for _, issue := range issues {
		go func(issue repo.Issue) {
			defer wg.Done()
			c := crawler.NewCrawler()
			cfg := issueConfigFromLabels(issue.Labels)

			posts := c.Collect(issue.Title, from, to)
			tr := crawler.NewTrending(posts)

			switch {
			case cfg.deviate > 0:
				r.CreateTrendingReport(issue.ID, tr.Deviate(cfg.deviate))
			}
		}(issue)
	}
	wg.Wait()

	return nil
}
