package crawler

import (
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Crawler) parseBoardPage(page string, from time.Time, to time.Time) ([]Post, string, bool) {
	infos, prev, cont := c.parsePostInfos(page, from, to)
	postc := make(chan Post)
	posts := make([]Post, 0)

	var wg sync.WaitGroup
	wg.Add(len(infos))
	for _, info := range infos {
		go func(info postInfo) {
			defer wg.Done()
			postc <- c.parsePost(info)
		}(info)
	}

	go func() {
		wg.Wait()
		close(postc)
	}()

	for post := range postc {
		posts = append(posts, post)
	}

	return posts, prev, cont
}

func (c *Crawler) parsePostInfos(page string, from time.Time, to time.Time) (map[string]postInfo, string, bool) {
	doc, err := c.loader.Load(page)
	if err != nil {
		panic(err)
	}

	prev := doc.Find(".btn-group-paging .btn").Eq(1).AttrOr("href", "")
	infos := make(map[string]postInfo)
	cont := true
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		Each(func(i int, sel *goquery.Selection) {
			title := sel.Find(".title > a")
			href, ok := title.Attr("href")
			if !ok {
				return
			}

			id, createAt := parseURL(href)
			if createAt.Before(from) {
				cont = false
				return
			} else if createAt.Before(to) {
				infos[id] = postInfo{
					URL:      c.resolver.getFullURL(href),
					CreateAt: createAt,
				}
			}
		})

	return infos, c.resolver.getFullURL(prev), cont
}

func (c *Crawler) parsePost(info postInfo) Post {
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
		NumPush:  push.Length(),
		NumUp:    push.FilterFunction(isPushUp).Length(),
		NumDown:  push.FilterFunction(isPushDown).Length(),
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
