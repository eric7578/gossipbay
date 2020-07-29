package crawler

import (
	"testing"
	"time"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCrawler_parseBoardPage(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	c := NewCrawler()
	c.loader = &testutil.TestDataLoader{}
	posts, _, ok := c.parsePostInfos("testdata/board.html", time.Date(2020, 7, 3, 0, 0, 0, 0, loc))

	assert.True(t, ok)
	assert.Equal(t, 19, len(posts))

	_, _, ok = c.parsePostInfos("testdata/board.html", time.Date(2020, 7, 5, 0, 0, 0, 0, loc))
	assert.False(t, ok)
}

func TestCrawler_parsePost(t *testing.T) {
	c := NewCrawler()
	c.loader = &testutil.TestDataLoader{}
	p := c.parsePost(postInfo{URL: "testdata/M.1593841729.A.BDA.html"})

	assert.Equal(t, "[閒聊] 聊聊大王", p.Title)
	assert.Equal(t, "sky419012 (fly)", p.Author)
	assert.Equal(t, 36, p.NumPush)
	assert.Equal(t, 14, p.NumUp)
	assert.Equal(t, 5, p.NumDown)
}
