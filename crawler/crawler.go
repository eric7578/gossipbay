package crawler

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DocumentLoader interface {
	Load(string) (*goquery.Document, error)
}

type Crawler struct {
	loader DocumentLoader
}

func NewCrawler() *Crawler {
	return &Crawler{
		loader: &HttpLoader{},
	}
}

func (c *Crawler) CollectUntil(board string, until time.Time) []Post {
	var pagePosts []Post
	posts := make([]Post, 0)
	cont := true
	page := "/bbs/" + board + "/index.html"
	for {
		pagePosts, page, cont = c.parseBoardPage(page, until)
		posts = append(posts, pagePosts...)
		if !cont {
			break
		}
	}
	return posts
}
