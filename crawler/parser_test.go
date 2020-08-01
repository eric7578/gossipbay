package crawler

import (
	"testing"
	"time"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_pageParser_ParsePostList(t *testing.T) {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Taipei")
	p := pageParser{
		ldr: &testutil.TestDataLoader{},
	}
	posts, _, ok := p.ParsePostList("testdata/board.html", time.Date(2020, 7, 3, 0, 0, 0, 0, loc), now)

	assert.True(t, ok)
	assert.Equal(t, 19, len(posts))

	_, _, ok = p.ParsePostList("testdata/board.html", time.Date(2020, 7, 5, 0, 0, 0, 0, loc), now)
	assert.False(t, ok)
}

func Test_pageParser_ParsePost(t *testing.T) {
	p := pageParser{
		ldr: &testutil.TestDataLoader{},
	}
	post := p.ParsePost("testdata/M.1593841729.A.BDA.html")

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
