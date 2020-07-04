package crawler

import (
	"sort"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Crawler) parseBoardPage(infoc chan PostInfo, page string, until time.Time) (prev string, next bool) {
	doc, err := c.loader.Load(page)
	if err != nil {
		panic(err)
	}

	// nav btns
	btns := doc.Find(".btn-group-paging .btn")
	prev = btns.Eq(1).AttrOr("href", "")

	var wg sync.WaitGroup
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		EachWithBreak(func(i int, sel *goquery.Selection) bool {
			rowInfo := parsePostInfo(sel)
			if rowInfo.CreateAt.Before(until) {
				next = false
				return false
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				url := sel.Find(".dropdown > .item > a").Eq(0).AttrOr("href", "")
				if info, ok := c.getSameTitledPostInfos(url); ok {
					infoc <- info
				} else {
					// in some edge cases, you may not find any same titled posts
					// which might be caused by the name of the title
					infoc <- parsePostInfo(sel)
				}
			}()
			return true
		})

	wg.Wait()

	return
}

func (c *Crawler) getSameTitledPostInfos(sameTitleSearchPage string) (info PostInfo, ok bool) {
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

	numInfos := len(infos)
	switch numInfos {
	case 0:
		ok = false
	case 1:
		ok = true
		info = infos[0]
	default:
		ok = true
		sort.Sort(AscCreateDate(infos))
		info = infos[0]
		info.Replies = infos[1:]
	}
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
