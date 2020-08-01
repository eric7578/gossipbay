package crawler

import (
	"fmt"
	"sync"
	"time"
)

type CollectOption struct {
	Board string
	From  time.Time
	To    time.Time
}

type Crawler struct {
	parser Parser
}

func NewCrawler() *Crawler {
	return &Crawler{
		parser: &pageParser{
			ldr:    &httpLoader{},
			domain: "https://ptt.cc/bbs",
		},
	}
}

func (c *Crawler) Collect(opt CollectOption) []Post {
	var (
		posts []Post
		page  = fmt.Sprintf("/bbs/%s/index.html", opt.Board)
		postc = make(chan Post)
		wg    sync.WaitGroup
	)
	for {
		infos, next, more := c.parser.ParsePostList(page, opt.From, opt.To)
		wg.Add(len(infos))
		for _, info := range infos {
			go func(info PostInfo) {
				defer wg.Done()
				postc <- c.parser.ParsePost(info.URL)
			}(info)
		}
		if more {
			page = next
		} else {
			break
		}
	}

	go func() {
		wg.Wait()
		close(postc)
	}()

	for post := range postc {
		posts = append(posts, post)
	}

	return posts
}
