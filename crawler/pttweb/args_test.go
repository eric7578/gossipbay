package pttweb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseArgs(t *testing.T) {
	args, err := parseArgs(map[string]string{
		"board":   "Gossiping",
		"deviate": "0.9",
		"period":  "daily",
		"timeout": "2",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Gossiping", args.board)
	assert.Equal(t, 0.9, args.deviate)
	assert.True(t, args.timeout > 0)
	assert.Equal(t, time.Duration(2)*time.Second, args.timeout)
}

func Test_parseArgs_withoutTimeout(t *testing.T) {
	args, err := parseArgs(map[string]string{
		"board":   "Gossiping",
		"deviate": "0.9",
		"period":  "daily",
	})

	assert.Nil(t, err)
	assert.Equal(t, "Gossiping", args.board)
	assert.True(t, args.timeout == 0)
}
