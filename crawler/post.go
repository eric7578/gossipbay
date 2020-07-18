package crawler

import (
	"path"
	"strconv"
	"strings"
	"time"
)

func parseURL(href string) (id string, createAt time.Time) {
	timestamp, err := strconv.ParseInt(strings.Split(href, ".")[1], 10, 54)
	if err != nil {
		panic(err)
	}
	_, id = path.Split(href)
	createAt = time.Unix(timestamp, 0)
	return
}

type postInfo struct {
	URL      string
	CreateAt time.Time
	IsReply  bool
}

type Post struct {
	ID       string
	URL      string
	CreateAt time.Time
	Title    string
	Author   string
	NumPush  int
	NumUp    int
	NumDown  int
}

type Thread struct {
	Title   string
	NumPush int
	NumUp   int
	NumDown int
	Posts   []Post
}
