package schedule

import (
	"testing"
	"time"

	"github.com/eric7578/gossipbay/repo"
	"github.com/stretchr/testify/assert"
)

func Test_parseRunOption(t *testing.T) {
	issue := repo.Issue{
		Title:  "Gossiping",
		Labels: []string{"trending-0.9", "timeout-2"},
	}
	opt := parseRunOption(issue)

	assert.Equal(t, "Gossiping", opt.Board)
	assert.Equal(t, 0.9, opt.Deviate)
	assert.True(t, opt.Timeout > 0)
	assert.Equal(t, time.Duration(2)*time.Second, opt.Timeout)
}

func Test_parseRunOption_withoutTimeout(t *testing.T) {
	issue := repo.Issue{
		Labels: []string{},
	}
	opt := parseRunOption(issue)

	assert.True(t, opt.Timeout == 0)
}
