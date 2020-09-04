package crawler

import (
	"testing"
	"time"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPageCrawlerVisitBoard(t *testing.T) {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Taipei")
	p := PageCrawler{
		Loader: &testutil.TestDataLoader{},
	}
	posts, next, _ := p.VisitBoard("testdata/board.html", time.Date(2020, 7, 3, 0, 0, 0, 0, loc), now)

	assert.NotEqual(t, "", next)
	assert.Equal(t, 19, len(posts))

	_, next, _ = p.VisitBoard("testdata/board.html", time.Date(2020, 7, 5, 0, 0, 0, 0, loc), now)
	assert.Equal(t, "", next)
}
