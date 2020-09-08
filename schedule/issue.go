package schedule

import (
	"sync"
	"time"

	"github.com/eric7578/gossipbay/repo"
)

type IssueReports struct {
	Updated time.Time
	Boards  []IssueBoardReport
}

type IssueBoardReport struct {
	BoardReport
	Issue int
}

type RunIssueOptions struct {
	Labels []string
}

func (s *Scheduler) RunIssues(r *repo.Github, opt RunIssueOptions) (IssueReports, error) {
	issues := r.ListIssues(opt.Labels...)
	breportc := make(chan IssueBoardReport)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(issues))
		for _, issue := range issues {
			go func(issue repo.Issue) {
				defer wg.Done()
				if opt := parseRunOption(issue); opt.isValid() {
					report := IssueBoardReport{
						Issue: issue.ID,
					}
					if r, err := s.Run(opt); err != nil {
						// TODO: error report
					} else {
						report.BoardReport = r
					}
					breportc <- report
				}
			}(issue)
		}
		wg.Wait()
		close(breportc)
	}()

	ireports := IssueReports{
		Updated: time.Now(),
	}
	for j := range breportc {
		ireports.Boards = append(ireports.Boards, j)
	}

	return ireports, nil
}
