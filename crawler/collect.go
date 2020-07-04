package crawler

import (
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Crawler) CollectUntil(page string, until time.Time) <-chan map[string]PostInfo {
	infoc := make(chan map[string]PostInfo)
	go func() {
		defer close(infoc)
		for {
			prev, cont := c.parseBoardPage(infoc, page, until)
			if cont {
				page = prev
			} else {
				return
			}
		}
	}()
	return infoc
}

func (c *Crawler) parseBoardPage(infoc chan map[string]PostInfo, page string, until time.Time) (prev string, cont bool) {
	doc, err := c.loader.Load(page)
	if err != nil {
		panic(err)
	}

	prev = doc.Find(".btn-group-paging .btn").Eq(1).AttrOr("href", "")
	infos, cont := parsePostInfos(doc, until)

	var wg sync.WaitGroup
	for _, info := range infos {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			doc, err := c.loader.Load(url)
			if err != nil {
				panic(err)
			}

			infos, _ := parsePostInfos(doc, until)
			infoc <- infos
		}(info.SameTitleURL)
	}
	wg.Wait()

	return
}

func parsePostInfos(doc *goquery.Document, until time.Time) (infos map[string]PostInfo, cont bool) {
	infos = make(map[string]PostInfo)
	cont = true
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		EachWithBreak(func(i int, sel *goquery.Selection) bool {
			title := sel.Find(".title > a")
			href, ok := title.Attr("href")
			if !ok {
				return true
			}

			id, createAt := parseURL(href)
			if createAt.Before(until) {
				cont = false
				return false
			}

			infos[id] = PostInfo{
				URL:          href,
				SameTitleURL: sel.Find(".dropdown > .item > a").Eq(0).AttrOr("href", ""),
				CreateAt:     createAt,
				Relates:      make([]PostInfo, 0),
			}
			return true
		})

	return
}

type AscCreateDate []PostInfo

func (a AscCreateDate) Len() int {
	return len(a)
}

func (a AscCreateDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a AscCreateDate) Less(i, j int) bool {
	return a[i].CreateAt.Before(a[j].CreateAt)
}
