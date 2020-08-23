package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eric7578/gossipbay/crawler"
)

type Github struct {
	token string
	repo  string
	owner string
}

func NewGithub(repository, token string) *Github {
	segs := strings.Split(repository, "/")
	return &Github{
		token: token,
		owner: segs[0],
		repo:  segs[1],
	}
}

func (gh *Github) ListIssues(labels ...string) []Issue {
	type GithubIssue struct {
		Number      int
		Title       string
		Labels      []struct{ Name string }
		PullRequest json.RawMessage `json:"pull_request,omitempty"`
	}
	githubIssues := make([]GithubIssue, 0)
	err := gh.api("GET", fmt.Sprintf("/repos/%s/%s/issues", gh.owner, gh.repo), &githubIssues, nil)
	if err != nil {
		panic(err)
	}

	issues := make([]Issue, 0)
	for _, githubIssue := range githubIssues {
		if len(githubIssue.PullRequest) == 0 {
			var labels []string
			for _, label := range githubIssue.Labels {
				labels = append(labels, label.Name)
			}
			issues = append(issues, Issue{
				ID:     githubIssue.Number,
				Title:  githubIssue.Title,
				Labels: labels,
			})
		}
	}
	return issues
}

func (gh *Github) CreateTrendingReport(issueID int, threads []crawler.Thread) {
	type CreateTrendingReportBody struct {
		Body string `json:"body"`
	}

	var buf bytes.Buffer
	if err := mdTmpl.Execute(&buf, threads); err != nil {
		panic(err)
	}
	payload := CreateTrendingReportBody{
		Body: buf.String(),
	}

	err := gh.api("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/comments", gh.owner, gh.repo, issueID), nil, &payload)
	if err != nil {
		panic(err)
	}
}

func (gh *Github) api(method string, path string, response interface{}, body interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	c := http.Client{}
	url := "https://api.github.com" + path
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		panic(err)
	}

	if bodyReader != nil {
		req.Header.Add("Accept", "application/vnd.github.v3+json")
	}

	if gh.token != "" {
		req.Header.Add("Authorization", "token "+gh.token)
	}

	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("request to %s got status code error: [%d] %s", url, res.StatusCode, res.Status)
	}

	defer res.Body.Close()
	if response != nil {
		json.NewDecoder(res.Body).Decode(response)
	}
	return nil
}
