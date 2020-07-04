package crawler

import (
	"testing"
	"time"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCrawler_parsePostInfos(t *testing.T) {
	ldr := testutil.TestDataLoader{}
	doc, _ := ldr.Load("testdata/board.html")
	loc, _ := time.LoadLocation("Asia/Taipei")
	posts, ok := parsePostInfos(doc, time.Date(2020, 7, 3, 0, 0, 0, 0, loc))

	assert.True(t, ok)
	assert.Equal(t, 19, len(posts))

	_, ok = parsePostInfos(doc, time.Date(2020, 7, 5, 0, 0, 0, 0, loc))
	assert.False(t, ok)
}
