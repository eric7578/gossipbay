package crawler

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DocumentLoader interface {
	Load(string) (*goquery.Document, error)
}

type Crawler struct {
	board  string
	loader DocumentLoader
}

func NewCrawler(board string) *Crawler {
	return &Crawler{
		board:  board,
		loader: &HttpLoader{},
	}
}

func (c *Crawler) CollectUntil(t *Trending, until time.Time) {
	var (
		posts []Post
		next  = true
		page  = "/bbs/" + c.board + "/index.html"
	)
	for next {
		posts, page, next = c.parseBoardPage(page, until)
		t.addPosts(posts...)
	}
}
