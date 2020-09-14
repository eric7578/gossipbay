package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseRunOption(t *testing.T) {
	issue := GithubIssue{
		Title: "Gossiping",
		Labels: []struct{ Name string }{
			{Name: "trending-0.9"},
			{Name: "timeout-2"},
		},
	}
	opt := parseTrendingOption(issue)

	assert.Equal(t, "Gossiping", opt.Board)
	assert.Equal(t, 0.9, opt.Deviate)
	assert.True(t, opt.Timeout > 0)
	assert.Equal(t, time.Duration(2)*time.Second, opt.Timeout)
}

func Test_parseRunOption_withoutTimeout(t *testing.T) {
	issue := GithubIssue{
		Labels: []struct{ Name string }{},
	}
	opt := parseTrendingOption(issue)

	assert.True(t, opt.Timeout == 0)
}
