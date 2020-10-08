package ptt

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var (
	regProtocol = regexp.MustCompile("https?://")
)

func getFullURL(s string) string {
	if regProtocol.MatchString(s) {
		return s
	}

	u, err := url.Parse("https://www.ptt.cc")
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, s)
	return u.String()
}

type Loader interface {
	Load(string) (*goquery.Document, error)
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

type PTTCrawler struct {
	Loader
}

func NewPTTCrawler() *PTTCrawler {
	return &PTTCrawler{
		Loader: &httpLoader{},
	}
}
