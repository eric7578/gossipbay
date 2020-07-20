package repo

type Label struct {
	Name string
}

type Issue struct {
	ID     int
	Title  string
	Labels []Label
}

type Repository interface {
	ListIssues() []Issue
	CreateIssueComment(issueID int, content string)
}
