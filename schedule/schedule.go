package schedule

import (
	"sync"

	"github.com/eric7578/gossipbay/crawler"
	"github.com/eric7578/gossipbay/repo"
)

func RunSchedule(envTrending string, r repo.Repository) {
	var wg sync.WaitGroup
	issues := r.ListIssues(envTrending)
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
}
