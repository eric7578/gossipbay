package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type HttpLoader struct {
	sem chan struct{}
}

func (ldr *HttpLoader) Load(p string) (doc *goquery.Document, err error) {
	ldr.sem <- struct{}{}
	defer func() {
		<-ldr.sem
	}()

	c := http.Client{}
	req, err := http.NewRequest("GET", "https://www.ptt.cc"+p, nil)
	req.Header.Add("cookie", "over18=1")

	res, err := c.Do(req)
	if err == nil {
		if res.StatusCode != 200 {
			err = fmt.Errorf("page: %s got status code error: [%d] %s", p, res.StatusCode, res.Status)
		} else {
			defer res.Body.Close()
			doc, err = goquery.NewDocumentFromReader(res.Body)
		}
	}
	return
}
