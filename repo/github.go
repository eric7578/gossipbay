package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GithubIssue struct {
	Title  string
	Labels []struct {
		Name string
	}
}

type Github struct {
	apiPath string
	token   string
}

func NewGithub(repository, token string) *Github {
	segs := strings.Split(repository, "/")
	return &Github{
		apiPath: fmt.Sprintf("/repos/%s/%s", segs[0], segs[1]),
		token:   token,
	}
}

func (gh *Github) GetIssues(labels ...string) ([]GithubIssue, error) {
	issues := []GithubIssue{}
	query := ""
	if len(labels) > 0 {
		query = "?labels=" + strings.Join(labels, ",")
	}
	err := gh.api("GET", gh.apiPath+"/issues"+query, &issues, nil)
	return issues, err
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
