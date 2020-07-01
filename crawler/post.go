package crawler

import (
	"time"
)

type PostInfo struct {
	ID       string
	Author   string
	Title    string
	CreateAt time.Time
	URL      string
	Replies  []PostInfo
}

type Post struct {
	Info    PostInfo
	Replies []Post
	NumPush int
	NumUp   int
	NumDown int
}
