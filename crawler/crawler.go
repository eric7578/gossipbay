package crawler

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DocumentLoader interface {
	Load(string) (*goquery.Document, error)
}

type Crawler struct {
	resolver *resolver
	loader   DocumentLoader
}

func NewCrawler() *Crawler {
	return &Crawler{
		resolver: &resolver{
			domain: "https://www.ptt.cc",
		},
		loader: &HttpLoader{},
	}
}

func (c *Crawler) CollectUntil(board string, t *Trending, until time.Time) {
	var (
		posts []Post
		next  = true
		page  = c.resolver.getBoardIndex(board)
	)
	for next {
		posts, page, next = c.parseBoardPage(page, until)
		t.addPosts(posts...)
	}
}
