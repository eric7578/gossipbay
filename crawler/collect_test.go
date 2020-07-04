package crawler

import (
	"testing"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCrawler_getSameTitledPostInfos(t *testing.T) {
	c := Crawler{
		loader: &testutil.TestDataLoader{},
	}
	info, _ := c.getSameTitledPostInfos("testdata/sametitle.html")

	assert.Equal(t, "/bbs/Gossiping/M.1593826452.A.02F.html", info.URL)

	assert.Equal(t, 2, len(info.Relates))
	assert.Equal(t, "/bbs/Gossiping/M.1593827377.A.6C2.html", info.Relates[0].URL)
	assert.Equal(t, "/bbs/Gossiping/M.1593827873.A.2F9.html", info.Relates[1].URL)
}
