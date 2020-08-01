package repo

import "github.com/eric7578/gossipbay/crawler"

type Issue struct {
	ID     int
	Title  string
	Labels []string
}

type Repository interface {
	ListIssues(...string) []Issue
	CreateTrendingReport(issueID int, threads []crawler.Thread)
}
