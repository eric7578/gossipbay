package repo

const (
	// period labels
	TrendingWeekly = "trending-weekly"
	TrendingDaily  = "trending-daily"
)

type Issue struct {
	ID     int
	Title  string
	Labels map[string]struct{}
}

type Repository interface {
	ListIssues(...string) []Issue
	CreateIssueComment(issueID int, content string)
}
