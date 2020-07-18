package crawler

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DocumentLoader interface {
	Load(string) (*goquery.Document, error)
}

type Crawler struct {
	loader DocumentLoader
}

func NewCrawler() *Crawler {
	return &Crawler{
		loader: &HttpLoader{},
	}
}

func (c *Crawler) CollectUntil(board string, until time.Time) []Thread {
	threads := make([]Thread, 0)
	titles := make(map[string]*Thread)
	var pagePosts []Post
	cont := true
	page := "/bbs/" + board + "/index.html"
	for {
		pagePosts, page, cont = c.parseBoardPage(page, until)
		for _, post := range pagePosts {
			title := trimTitle(post.Title)
			if t, ok := titles[title]; ok {
				t.NumPush += post.NumPush
				t.NumUp += post.NumUp
				t.NumDown += post.NumDown
				t.Posts = append(t.Posts, post)
			} else {
				t := Thread{
					Title:   post.Title,
					NumPush: post.NumPush,
					NumUp:   post.NumUp,
					NumDown: post.NumDown,
					Posts:   []Post{post},
				}
				threads = append(threads, t)
				titles[title] = &t
			}
		}
		if !cont {
			break
		}
	}
	return threads
}

var leadingReply = regexp.MustCompile(`^Re:\s*`)

func trimTitle(title string) string {
	title = strings.Trim(title, " ")
	return leadingReply.ReplaceAllString(title, "")
}
