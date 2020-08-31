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

// type Crawler struct {
// 	parser Parser
// }

// func NewCrawler() *Crawler {
// 	return &Crawler{
// 		parser: &PageParser{
// 			ldr:    &httpLoader{},
// 			domain: "https://www.ptt.cc",
// 		},
// 	}
// }

// func (c *Crawler) Collect(ctx context.Context, board string, from, to time.Time) []Post {
// 	var (
// 		posts []Post
// 		page  = fmt.Sprintf("/bbs/%s/index.html", board)
// 		postc = make(chan Post)
// 		wg    sync.WaitGroup
// 	)
// 	for {
// 		infos, next, more := c.parser.VisitPostList(page, from, to)
// 		wg.Add(len(infos))
// 		for _, info := range infos {
// 			go func(info PostInfo) {
// 				defer wg.Done()
// 				select {
// 				case postc <- c.parser.VisitPost(info.URL):
// 				case <-ctx.Done():
// 					return
// 				}
// 			}(info)
// 		}
// 		if more {
// 			page = next
// 		} else {
// 			break
// 		}
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(postc)
// 	}()

// 	for post := range postc {
// 		posts = append(posts, post)
// 	}

// 	return posts
// }
