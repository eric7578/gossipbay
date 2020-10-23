package pttweb

import (
	"context"
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

type PttWorker struct {
	Loader
}

func NewPttWorker() *PttWorker {
	return &PttWorker{
		Loader: &httpLoader{},
	}
}

func (w *PttWorker) Accept(args map[string]string) bool {
	return args["_type"] == "pttweb"
}

func (w *PttWorker) Run(args map[string]string) (interface{}, error) {
	targs, err := parseArgs(args)
	if err != nil {
		return nil, err
	}
	return w.trending(context.Background(), targs)
}
