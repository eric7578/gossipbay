package crawler

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	regProtocol = regexp.MustCompile("https?://")
)

type Loader interface {
	Load(string) (*goquery.Document, error)
}

type Crawler interface {
	VisitBoard(page string, from time.Time, to time.Time) ([]PostInfo, string, bool)
	VisitPost(page string) Post
}

type httpLoader struct {
}

func (ldr *httpLoader) Load(url string) (doc *goquery.Document, err error) {
	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("cookie", "over18=1")

	res, err := c.Do(req)
	if err == nil {
		if res.StatusCode != 200 {
			err = fmt.Errorf("page: %s got status code error: [%d] %s", url, res.StatusCode, res.Status)
		} else {
			defer res.Body.Close()
			doc, err = goquery.NewDocumentFromReader(res.Body)
		}
	}
	return
}

type PageCrawler struct {
	Loader
	domain string
}

func NewPageCrawler() *PageCrawler {
	return &PageCrawler{
		Loader: &httpLoader{},
		domain: "https://www.ptt.cc",
	}
}
