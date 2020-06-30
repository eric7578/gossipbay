package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type HttpLoader struct {
}

func (ldr *HttpLoader) Load(p string) (doc *goquery.Document, err error) {
	c := http.Client{}
	req, err := http.NewRequest("GET", "https://www.ptt.cc"+p, nil)
	req.Header.Add("cookie", "over18=1")

	res, err := c.Do(req)
	if err == nil {
		if res.StatusCode != 200 {
			err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		} else {
			defer res.Body.Close()
			doc, err = goquery.NewDocumentFromReader(res.Body)
		}
	}
	return
}
