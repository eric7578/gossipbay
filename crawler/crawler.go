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

	urls := doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		Map(func(i int, sel *goquery.Selection) string {
			return sel.Find(".dropdown > .item > a").Eq(0).AttrOr("href", "")
		})

	var wg sync.WaitGroup
	wg.Add(len(urls))
	infoc := make(chan []PostInfo)

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			infoc <- c.getSameTitledPostInfos(url)
		}(url)
	}

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
			// EX: https://www.ptt.cc/bbs/Gossiping/M.1592706173.A.56E.html
			title := sel.Find(".title > a")
			href, _ := title.Attr("href")
			_, file := path.Split(href)
			createAt, err := strconv.ParseInt(strings.Split(file, ".")[1], 10, 54)
			if err != nil {
				panic(err)
			}

			info := PostInfo{
				ID:       file,
				Author:   sel.Find(".author").Text(),
				Title:    title.Text(),
				CreateAt: time.Unix(createAt, 0),
				URL:      href,
				Replies:  make([]PostInfo, 0),
			}

			infos = append(infos, info)
		})

	return infos
}
