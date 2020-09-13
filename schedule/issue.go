package schedule

import (
	"sync"

	"github.com/eric7578/gossipbay/repo"
)

type RunIssueOptions struct {
	Labels []string
}

type BoardTrending struct {
	Board   string    `json:"board"`
	Threads []*Thread `json:"threads"`
}

func (s *Scheduler) RunIssues(r *repo.Github, opt RunIssueOptions) ([]BoardTrending, error) {
	issues := r.ListIssues(opt.Labels...)
	boardc := make(chan BoardTrending)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(issues))
		for _, issue := range issues {
			go func(issue repo.Issue) {
				defer wg.Done()
				if opt := parseTrendingOption(issue); opt.isValid() {
					if threads, err := s.Trending(opt); err != nil {
						// TODO: error report
					} else {
						boardc <- BoardTrending{
							Board:   opt.Board,
							Threads: threads,
						}
					}
				}
			}(issue)
		}
		wg.Wait()
		close(boardc)
	}()

	boards := make([]BoardTrending, 0)
	for b := range boardc {
		boards = append(boards, b)
	}

	return boards, nil
}
