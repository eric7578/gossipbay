package crawler

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (p *PageCrawler) VisitPost(page string) (Post, error) {
	post := Post{}
	doc, err := p.Load(page)
	if err != nil {
		return post, err
	}

	id, createdAt := parseURL(page)
	post.ID = id
	post.URL = page
	post.CreatedAt = createdAt
	post.parseMeta(doc.Find(".article-meta-tag"))
	post.parsePush(doc.Find(".push"))
	post.parseContent(doc.Find("#main-content"))
	return post, nil
}

type Post struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Re          bool      `json:"reply"`
	Author      string    `json:"author"`
	TextContent string    `json:"textContent"`
	PushTotal   int       `json:"pushTotal"`
	BattlePush  int       `json:"battlePush"`
	UniquePush  int       `json:"uniquePush"`
	PushUp      int       `json:"pushUp"`
	PushDown    int       `json:"pushDown"`
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

func (p *Post) parseContent(mainContainer *goquery.Selection) {
	clone := mainContainer.Clone()
	clone.Children().Each(func(i int, sel *goquery.Selection) {
		// TODO: keep highlighted text nodes
		// Ref: https://www.ptt.cc/bbs/Grasshopper/M.1246471695.A.FCA.html
		sel.Remove()
	})
	p.TextContent = strings.ReplaceAll(strings.TrimSpace(clone.Text()), "\n", "")
}

func (p *Post) parseMeta(tags *goquery.Selection) {
	tags.Each(func(i int, sel *goquery.Selection) {
		tag := sel.Text()
		val := sel.SiblingsFiltered(".article-meta-value").Text()
		switch tag {
		case "標題":
			p.Title = val
			p.Re = strings.Index(val, "Re:") == 0
		case "作者":
			p.Author = val
		}
	})
}

func (p *Post) parsePush(pushLines *goquery.Selection) {
	p.PushTotal = pushLines.Length()

	authors := make(map[string]struct{})
	lastPushTime := make(map[string]time.Time)

	pushLines.Each(func(i int, sel *goquery.Selection) {
		pushTag := strings.TrimSpace(sel.Find(".push-tag").Text())
		pushedAt, _ := time.Parse("01/01 15:04", strings.TrimSpace(sel.Find(".push-ipdatetime").Text()))
		up := pushTag == "推"
		down := pushTag == "噓"
		author := strings.TrimSpace(sel.Find(".push-userid").Text())

		if up {
			p.PushUp += 1
		} else if down {
			p.PushDown += 1
		}

		if _, ok := authors[author]; !ok {
			p.UniquePush += 1
			authors[author] = struct{}{}
		}

		prevPushedAt, ok := lastPushTime[author]
		if !ok || prevPushedAt.Before(pushedAt.Add(-3*time.Minute)) {
			p.BattlePush += 1
			lastPushTime[author] = pushedAt
		}
	})
}
