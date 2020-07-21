package crawler

import (
	"regexp"
	"strings"
)

var leadingReply = regexp.MustCompile(`^Re:\s*`)

func trimTitle(title string) string {
	title = strings.Trim(title, " ")
	return leadingReply.ReplaceAllString(title, "")
}

type Thread struct {
	Title   string
	NumPush int
	NumUp   int
	NumDown int
	Posts   []Post
}

type Trending struct {
	threads []Thread
	titles  map[string]*Thread
}

func NewTrending() *Trending {
	return &Trending{
		titles: make(map[string]*Thread),
	}
}

func (tr *Trending) addPosts(posts ...Post) {
	for _, post := range posts {
		title := trimTitle(post.Title)
		t, ok := tr.titles[title]
		if ok {
			t.NumPush += post.NumPush
			t.NumUp += post.NumUp
			t.NumDown += post.NumDown
			t.Posts = append(t.Posts, post)
		} else {
			t = &Thread{
				Title:   post.Title,
				NumPush: post.NumPush,
				NumUp:   post.NumUp,
				NumDown: post.NumDown,
				Posts:   []Post{post},
			}
			tr.threads = append(tr.threads, *t)
			tr.titles[title] = t
		}
	}
}

func (tr *Trending) Deviate(v float32) []Thread {
	var (
		numMaxPush float32
		numMinPush float32
	)
	for i, t := range tr.threads {
		npush := float32(t.NumPush)
		if i == 0 {
			numMaxPush = npush
			numMinPush = npush
		} else {
			if numMaxPush < npush {
				numMaxPush = npush
			}
			if numMinPush > npush {
				numMinPush = npush
			}
		}
	}

	var (
		top     = 1 - v
		offset  = numMaxPush - numMinPush
		threads []Thread
	)
	for _, t := range tr.threads {
		if float32(t.NumPush)/offset > top {
			threads = append(threads, t)
		}
	}
	return threads
}
