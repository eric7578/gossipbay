package crawler

import (
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	regProtocol = regexp.MustCompile("https?://")
)

type Parser interface {
	ParsePostList(page string, from time.Time, to time.Time) ([]PostInfo, string, bool)
	ParsePost(page string) Post
}

type pageParser struct {
	ldr    Loader
	domain string
}

func (p *pageParser) ParsePostList(page string, from time.Time, to time.Time) ([]PostInfo, string, bool) {
	doc, err := p.ldr.Load(p.getFullURL(page))
	if err != nil {
		panic(err)
	}

	var infos []PostInfo
	prev := doc.Find(".btn-group-paging .btn").Eq(1).AttrOr("href", "")
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

			_, createAt := parseURL(href)
			if createAt.Before(from) {
				cont = false
				return
			} else if createAt.Before(to) {
				infos = append(infos, PostInfo{
					URL:      p.getFullURL(href),
					CreateAt: createAt,
				})
			}
		})

	return infos, p.getFullURL(prev), cont
}

func (p *pageParser) ParsePost(page string) Post {
	doc, err := p.ldr.Load(page)
	if err != nil {
		panic(err)
	}

	id, createAt := parseURL(page)
	pushes := doc.Find(".push")
	metaTags := doc.Find(".article-meta-tag")
	noRepeatPush := set{}
	noRepeatUp := set{}
	noRepeatDown := set{}

	post := Post{
		ID:       id,
		URL:      page,
		CreateAt: createAt,
		Title:    metaTags.FilterFunction(isTitleMetaTag).Next().Text(),
		Author:   metaTags.FilterFunction(isAuthorMetaTag).Next().Text(),
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

	return post
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

func (p *pageParser) getFullURL(s string) string {
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
