package crawler

import (
	"strings"
	"sync"

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

func (c *Crawler) ParsePost(info PostInfo) Post {
	doc, err := c.loader.Load(info.URL)
	if err != nil {
		panic(err)
	}

	id, createAt := parseURL(info.URL)
	push := doc.Find(".push")
	metaTags := doc.Find(".article-meta-tag")

	p := Post{
		ID:       id,
		URL:      info.URL,
		CreateAt: createAt,
		Title:    metaTags.FilterFunction(isTitleMetaTag).Next().Text(),
		Author:   metaTags.FilterFunction(isAuthorMetaTag).Next().Text(),
		Replies:  make([]Post, 0),
		NumPush:  push.Length(),
		NumUp:    push.FilterFunction(isPushUp).Length(),
		NumDown:  push.FilterFunction(isPushDown).Length(),
	}

	if len(info.Relates) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(info.Relates))
		replyc := make(chan Post)

		for _, info := range info.Relates {
			go func(info PostInfo) {
				defer wg.Done()
				replyc <- c.ParsePost(info)
			}(info)
		}

		go func() {
			wg.Wait()
			close(replyc)
		}()

		for reply := range replyc {
			p.Replies = append(p.Replies, reply)
		}
	}

	return p
}

func isTitleMetaTag(i int, sel *goquery.Selection) bool {
	return strings.TrimSpace(sel.Text()) == "標題"
}

func isAuthorMetaTag(i int, sel *goquery.Selection) bool {
	return strings.TrimSpace(sel.Text()) == "作者"
}

func isPushUp(i int, sel *goquery.Selection) bool {
	return strings.TrimSpace(sel.Find(".push-tag").Text()) == "推"
}

func isPushDown(i int, sel *goquery.Selection) bool {
	return strings.TrimSpace(sel.Find(".push-tag").Text()) == "噓"
}
