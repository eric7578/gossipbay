package crawler

import (
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DocumentLoader interface {
	Load(string) (*goquery.Document, error)
}

type Crawler struct {
	loader DocumentLoader
}

func NewCrawler(ldr DocumentLoader) *Crawler {
	return &Crawler{
		loader: ldr,
	}
}

func (c *Crawler) ParsePostInfos(page string) (infos map[string]PostInfo, prev string, next string) {
	doc, err := c.loader.Load(page)
	if err != nil {
		panic(err)
	}

	// nav btns
	btns := doc.Find(".btn-group-paging .btn")
	prev = btns.Eq(1).AttrOr("href", "")
	next = btns.Eq(2).AttrOr("href", "")

	var wg sync.WaitGroup
	infoc := make(chan []PostInfo)
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		Each(func(i int, sel *goquery.Selection) {
			wg.Add(1)
			go func(sel *goquery.Selection) {
				defer wg.Done()
				url := sel.Find(".dropdown > .item > a").Eq(0).AttrOr("href", "")
				if infos := c.getSameTitledPostInfos(url); len(infos) != 0 {
					infoc <- infos
				} else {
					// in some edge cases, you may not find any same titled posts
					// which might be caused by the name of the title
					infoc <- []PostInfo{parsePostInfo(sel)}
				}
			}(sel)
		})

	go func() {
		wg.Wait()
		close(infoc)
	}()

	infos = make(map[string]PostInfo)
	for sameTitleInfos := range infoc {
		lastIndex := len(sameTitleInfos) - 1
		origin := sameTitleInfos[lastIndex]
		origin.Replies = sameTitleInfos[:lastIndex]
		infos[origin.ID] = origin
	}

	return
}

func (c *Crawler) getSameTitledPostInfos(sameTitleSearchPage string) []PostInfo {
	doc, err := c.loader.Load(sameTitleSearchPage)
	if err != nil {
		panic(err)
	}

	infos := make([]PostInfo, 0)
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		Each(func(i int, sel *goquery.Selection) {
			infos = append(infos, parsePostInfo(sel))
		})

	return infos
}

func (c *Crawler) ParsePost(info PostInfo) (p Post) {
	doc, err := c.loader.Load(info.URL)
	if err != nil {
		panic(err)
	}

	push := doc.Find(".push")

	p.Info = info
	p.NumPush = push.Length()
	p.NumUp = push.FilterFunction(isPushUp).Length()
	p.NumDown = push.FilterFunction(isPushDown).Length()
	p.Replies = make([]Post, 0)

	if len(info.Replies) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(info.Replies))
		replyc := make(chan Post)

		for _, info := range info.Replies {
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

	return
}

func parsePostInfo(sel *goquery.Selection) PostInfo {
	// EX: https://www.ptt.cc/bbs/Gossiping/M.1592706173.A.56E.html
	title := sel.Find(".title > a")
	href, _ := title.Attr("href")
	_, file := path.Split(href)
	createAt, err := strconv.ParseInt(strings.Split(file, ".")[1], 10, 54)
	if err != nil {
		panic(err)
	}

	return PostInfo{
		ID:       file,
		Author:   sel.Find(".author").Text(),
		Title:    title.Text(),
		CreateAt: time.Unix(createAt, 0),
		URL:      href,
		Replies:  make([]PostInfo, 0),
	}
}

func isPushUp(i int, sel *goquery.Selection) bool {
	return strings.TrimSpace(sel.Find(".push-tag").Text()) == "推"
}

func isPushDown(i int, sel *goquery.Selection) bool {
	return strings.TrimSpace(sel.Find(".push-tag").Text()) == "噓"
}
