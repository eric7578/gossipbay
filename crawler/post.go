package crawler

import (
	"time"
)

type PostInfo struct {
	URL      string
	CreateAt time.Time
	IsReply  bool
}

type Post struct {
	ID              string
	URL             string
	CreateAt        time.Time
	Title           string
	Author          string
	NumPush         int
	NumUp           int
	NumDown         int
	NumNoRepeatPush int
	NumNoRepeatUp   int
	NumNoRepeatDown int
}
