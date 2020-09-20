package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

type GithubIssue struct {
	Title  string
	Labels []struct{ Name string }
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

func (gh *Github) GetTrendingOptions(labels ...string) []crawler.TrendingOption {
	issues := make([]GithubIssue, 0)

	var query string
	if len(labels) > 0 {
		query = "?labels=" + strings.Join(labels, ",")
	}
	err := gh.api("GET", gh.apiPath+"/issues"+query, &issues, nil)
	if err != nil {
		panic(err)
	}

	opts := make([]crawler.TrendingOption, 0)
	for _, issue := range issues {
		opt := parseTrendingOption(issue)
		if opt.IsValid() {
			opts = append(opts, opt)
		}
	}
	return opts
}

func (gh *Github) PruneArtifact(daysAgo int) error {
	type GithubArtifacts struct {
		Artifacts []struct {
			ID        int `json:"id"`
			Expired   bool
			CreatedAt time.Time `json:"created_at"`
		}
	}

	var artifacts GithubArtifacts
	if err := gh.api("GET", gh.apiPath+"/actions/artifacts", &artifacts, nil); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errc := make(chan error)
	deadline := time.Now().Add(time.Duration(-24*daysAgo) * time.Hour)
	for _, artifact := range artifacts.Artifacts {
		if !artifact.Expired && artifact.CreatedAt.Before(deadline) {
			wg.Add(1)
			go func(artifactId int) {
				defer wg.Done()
				errc <- gh.api("DELETE", gh.apiPath+"/actions/artifacts/"+strconv.Itoa(artifactId), nil, nil)
			}(artifact.ID)
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
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
