package crawler

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Crawler) CollectUntil(page string, until time.Time) <-chan PostInfo {
	infoc := make(chan PostInfo)
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

func (c *Crawler) parseBoardPage(infoc chan PostInfo, page string, until time.Time) (prev string, cont bool) {
	doc, err := c.loader.Load(page)
	if err != nil {
		panic(err)
	}

	// nav btns
	btns := doc.Find(".btn-group-paging .btn")
	prev = btns.Eq(1).AttrOr("href", "")

	var wg sync.WaitGroup
	cont = true
	doc.
		Find(".r-list-container").
		Children().
		Filter(".search-bar").
		NextUntil(".r-list-sep").
		EachWithBreak(func(i int, sel *goquery.Selection) bool {
			rowInfo, err := parsePostInfo(sel)
			if errors.Is(err, errEmptyTitle) {
				return true
			}

			if rowInfo.CreateAt.After(until) {
				cont = false
				return false
			}

			if url, ok := sel.Find(".dropdown > .item > a").Eq(0).Attr("href"); ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					info, ok := c.getSameTitledPostInfos(url)
					if !ok {
						// in some edge cases, you may not find any same titled posts
						// which might be caused by the name of the title
						info, _ = parsePostInfo(sel)
					}
					infoc <- info
				}()
			}
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
			if info, err := parsePostInfo(sel); err == nil {
				infos = append(infos, info)
			}
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
		info.Relates = infos[1:]
	}
	return
}

func parsePostInfo(sel *goquery.Selection) (info PostInfo, err error) {
	title := sel.Find(".title > a")
	href, ok := title.Attr("href")
	if !ok {
		err = errEmptyTitle
		return
	}

	_, createAt := parseURL(href)
	info = PostInfo{
		CreateAt: createAt,
		URL:      href,
		Relates:  make([]PostInfo, 0),
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
