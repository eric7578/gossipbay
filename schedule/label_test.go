package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newScheduleOption(t *testing.T) {
	labels := []string{"trending-0.8"}
	_, _, deviate := parseIssueLabels(labels)

	assert.Equal(t, 0.8, deviate)
}
