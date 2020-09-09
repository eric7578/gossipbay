package crawler

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	ID              string    `json:"id"`
	URL             string    `json:"url"`
	CreatedAt       time.Time `json:"createdAt"`
	Title           string    `json:"title"`
	Re              bool      `json:"reply"`
	Author          string    `json:"author"`
	NumPush         int       `json:"numPush"`
	NumUp           int       `json:"numUp"`
	NumDown         int       `json:"numDown"`
	TextContent     string    `json:"textContent"`
	ExternalLinks   []string  `json:"exteralLinks"`
	NumNoRepeatPush int       `json:"-"`
	NumNoRepeatUp   int       `json:"-"`
	NumNoRepeatDown int       `json:"-"`
}

func (p *PageCrawler) VisitPost(page string) (post Post, err error) {
	doc, err := p.Load(page)
	if err != nil {
		return post, err
	}

	id, createAt := parseURL(page)
	pushes := doc.Find(".push")
	metaTags := doc.Find(".article-meta-tag")
	title := metaTags.FilterFunction(isTitleMetaTag).Next().Text()
	noRepeatPush := set{}
	noRepeatUp := set{}
	noRepeatDown := set{}

	post = Post{
		ID:        id,
		URL:       page,
		CreatedAt: createAt,
		Title:     title,
		Re:        strings.Index(title, "Re:") == 0,
		Author:    metaTags.FilterFunction(isAuthorMetaTag).Next().Text(),
	}

	pushes.Each(func(i int, sel *goquery.Selection) {
		pushTag := strings.TrimSpace(sel.Find(".push-tag").Text())
		author := strings.TrimSpace(sel.Find(".push-userid").Text())
		up := pushTag == "推"
		down := pushTag == "噓"
		noRepeatPush.add(author)
		if up {
			post.NumUp += 1
			noRepeatUp.add(author)
		} else if down {
			post.NumDown += 1
			noRepeatDown.add(author)
		}
	})

	post.NumPush = pushes.Length()
	post.NumNoRepeatPush = noRepeatPush.size()
	post.NumNoRepeatDown = noRepeatDown.size()
	post.NumNoRepeatUp = noRepeatUp.size()

	return post, nil
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

func parseURL(href string) (id string, createAt time.Time) {
	timestamp, err := strconv.ParseInt(strings.Split(path.Base(href), ".")[1], 10, 54)
	if err != nil {
		panic(err)
	}
	_, id = path.Split(href)
	createAt = time.Unix(timestamp, 0)
	return
}

type set map[string]struct{}

func (st set) add(s string) bool {
	_, ok := st[s]
	if !ok {
		st[s] = struct{}{}
	}
	return ok
}

func (st set) size() int {
	return len(st)
}
