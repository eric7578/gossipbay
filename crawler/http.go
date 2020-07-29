package crawler

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type HttpLoader struct {
}

func (ldr *HttpLoader) Load(url string) (doc *goquery.Document, err error) {
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
