package crawler

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CollectOption struct {
	Board string
	From  time.Time
	To    time.Time
}

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

func (c *Crawler) Collect(t *Trending, opt CollectOption) {
	if opt.To.IsZero() {
		opt.To = time.Now()
	}
	var (
		posts []Post
		next  = true
		page  = c.resolver.getBoardIndex(opt.Board)
	)
	for next {
		posts, page, next = c.parseBoardPage(page, opt.From, opt.To)
		t.addPosts(posts...)
	}
}
