package repo

import (
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

type Issue struct {
	ID     int
	Title  string
	Labels []string
}

type Comment struct {
	ID        int
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	ListIssues(...string) []Issue
	CreateTrendingReport(issueID int, threads []crawler.Thread)
	ListComments(since time.Time) []Comment
	RemoveComment(commentID int)
}
