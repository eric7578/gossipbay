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

type PostInfo struct {
	URL          string
	SameTitleURL string
	CreateAt     time.Time
	Relates      []PostInfo
}

type Post struct {
	ID       string
	URL      string
	CreateAt time.Time
	Title    string
	Author   string
	Replies  []Post
	NumPush  int
	NumUp    int
	NumDown  int
}
