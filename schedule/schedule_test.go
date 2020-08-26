package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_issueConfigFromLabels(t *testing.T) {
	labels := []string{"trending-0.8"}
	cfg := issueConfigFromLabels(labels)

	assert.Equal(t, 0.8, cfg.deviate)
}
