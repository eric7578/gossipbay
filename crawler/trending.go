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

type Trending map[string]*Thread

func NewTrending(posts []Post) Trending {
	threads := make(map[string]*Thread)

	for _, post := range posts {
		title := trimTitle(post.Title)
		if t, ok := threads[title]; ok {
			t.NumPush += post.NumPush
			t.NumUp += post.NumUp
			t.NumDown += post.NumDown
			t.Posts = append(t.Posts, post)
		} else {
			threads[title] = &Thread{
				Title:   post.Title,
				NumPush: post.NumPush,
				NumUp:   post.NumUp,
				NumDown: post.NumDown,
				Posts:   []Post{post},
			}
		}
	}
	return Trending(threads)
}

func (tr Trending) Deviate(v float32) []Thread {
	var (
		numMaxPush float32
		numMinPush float32
	)
	for _, t := range tr {
		npush := float32(t.NumPush)
		if numMaxPush < npush {
			numMaxPush = npush
		}
		if numMinPush > npush {
			numMinPush = npush
		}
	}

	var (
		top     = 1 - v
		offset  = numMaxPush - numMinPush
		threads []Thread
	)
	for _, t := range tr {
		if float32(t.NumPush)/offset > top {
			threads = append(threads, *t)
		}
	}
	return threads
}
