package crawler

import (
	"testing"
	"time"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/issue9/assert"
)

func TestPageCrawlerVisitPost(t *testing.T) {
	p := PageCrawler{
		Loader: &testutil.TestDataLoader{},
	}
	post, _ := p.VisitPost("testdata/M.1593841729.A.BDA.html")

	assert.Equal(t, "[閒聊] 聊聊大王", post.Title)
	assert.Equal(t, "sky419012 (fly)", post.Author)
	assert.Equal(t, 36, post.NumPush)
	assert.Equal(t, 14, post.NumUp)
	assert.Equal(t, 5, post.NumDown)
	assert.Equal(t, 14, post.NumNoRepeatPush)
	assert.Equal(t, 10, post.NumNoRepeatUp)
	assert.Equal(t, 3, post.NumNoRepeatDown)
}

func Test_parseURL(t *testing.T) {
	href := "https://ptt.cc/bbs/Gossiping/M.1592706173.A.56E.html"
	id, createAt := parseURL(href)

	assert.Equal(t, "M.1592706173.A.56E.html", id)
	loc, _ := time.LoadLocation("Asia/Taipei")
	assert.True(t, time.Date(2020, 6, 21, 10, 22, 53, 0, loc).Equal(createAt))
}
