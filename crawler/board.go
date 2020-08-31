package crawler

import (
	"net/url"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PostInfo struct {
	URL      string
	CreateAt time.Time
	IsReply  bool
}

func (p *PageCrawler) VisitBoard(page string, from time.Time, to time.Time) (infos []PostInfo, nextPage string, err error) {
	var doc *goquery.Document

	doc, err = p.Load(p.getFullURL(page))
	if err != nil {
		return infos, "", err
	}

	next := true
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

			_, createAt := parseURL(href)
			if createAt.Before(from) {
				next = false
				return
			} else if createAt.Before(to) {
				infos = append(infos, PostInfo{
					URL:      p.getFullURL(href),
					CreateAt: createAt,
				})
			}
		})

	if next {
		nextHref := doc.Find(".btn-group-paging .btn").Eq(1).AttrOr("href", "")
		nextPage = p.getFullURL(nextHref)
		return infos, nextPage, nil
	}

	return infos, "", nil
}

func (p *PageCrawler) getFullURL(s string) string {
	if regProtocol.MatchString(s) {
		return s
	}

	u, err := url.Parse(p.domain)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, s)
	return u.String()
}
