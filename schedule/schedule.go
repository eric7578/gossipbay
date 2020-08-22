package schedule

import (
	"fmt"
	"sync"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/eric7578/gossipbay/repo"
)

func RunSchedule(schedule string, r repo.Repository) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("schedule failed %s", e)
		}
	}()

	var wg sync.WaitGroup
	issues := r.ListIssues(schedule)
	wg.Add(len(issues))
	for _, issue := range issues {
		go func(issue repo.Issue) {
			defer wg.Done()
			c := crawler.NewCrawler()
			opt := newScheduleOption(issue.Labels)

			posts := c.Collect(issue.Title, opt.from, opt.to)
			tr := crawler.NewTrending(posts)

			switch {
			case opt.deviate > 0:
				r.CreateTrendingReport(issue.ID, tr.Deviate(opt.deviate))
			}
		}(issue)
	}
	wg.Wait()

	return nil
}
