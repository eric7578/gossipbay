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
