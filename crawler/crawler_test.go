package crawler

import (
	"testing"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCrawler_ParsePost(t *testing.T) {
	c := Crawler{
		loader: &testutil.TestDataLoader{},
	}
	p := c.ParsePost(PostInfo{URL: "testdata/M.1593841729.A.BDA.html"})

	assert.Equal(t, "[閒聊] 聊聊大王", p.Title)
	assert.Equal(t, "sky419012 (fly)", p.Author)
	assert.Equal(t, 36, p.NumPush)
	assert.Equal(t, 14, p.NumUp)
	assert.Equal(t, 5, p.NumDown)
}
