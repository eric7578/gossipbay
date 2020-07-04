package crawler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseURL(t *testing.T) {
	href := "/bbs/Gossiping/M.1592706173.A.56E.html"
	id, createAt := parseURL(href)

	assert.Equal(t, "M.1592706173.A.56E.html", id)
	loc, _ := time.LoadLocation("Asia/Taipei")
	assert.True(t, time.Date(2020, 6, 21, 10, 22, 53, 0, loc).Equal(createAt))
}
