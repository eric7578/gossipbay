package pttweb

import (
	"context"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type trendingArgs struct {
	board   string
	from    time.Time
	to      time.Time
	timeout time.Duration
	deviate float64
}

type Trending struct {
	Board   string    `json:"board"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
	Threads []Thread  `json:"threads"`
}

type Thread struct {
	Score   float64 `json:"score"`
	Deviate float64 `json:"deviate"`
	Posts   []Post  `json:"posts"`
}

func ScoreByBattle(p Post) float64 {
	return float64(p.BattlePush)
}

type trending struct {
	scoreFunc func(p Post) float64
	threads   map[string]*Thread
}

func newTrending(scoreFunc func(p Post) float64) *trending {
	return &trending{
		scoreFunc: scoreFunc,
		threads:   make(map[string]*Thread),
	}
}

func (t *trending) addPost(p Post) {
	group := genGroup(p.Title)
	thread, ok := t.threads[group]
	if ok {
		thread.Score += t.scoreFunc(p)
		thread.Posts = append(thread.Posts, p)
	} else {
		t.threads[group] = &Thread{
			Score: t.scoreFunc(p),
			Posts: []Post{p},
		}
	}
}

func (t *trending) deviate(threshold float64) []Thread {
	maxScore := float64(-1.0)
	minScore := float64(-1.0)

	for _, thread := range t.threads {
		if maxScore < thread.Score {
			maxScore = thread.Score
		}
		if minScore < 0 || minScore > thread.Score {
			minScore = thread.Score
		}
	}

	threads := make([]Thread, 0)
	offset := maxScore - minScore
	for _, thread := range t.threads {
		thread.Deviate = (thread.Score - minScore) / offset
		if threshold <= thread.Deviate {
			threads = append(threads, *thread)
		}
	}
	return threads
}

func genGroup(s string) string {
	reg := regexp.MustCompile(`^Re:\s*`)
	s = reg.ReplaceAllString(strings.TrimSpace(s), "")
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (w *PttWorker) trending(ctx context.Context, opt trendingArgs) (Trending, error) {
	var (
		cancel   context.CancelFunc
		trending = Trending{
			Board: opt.board,
			From:  opt.from,
			To:    opt.to,
		}
	)
	if opt.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	t := newTrending(ScoreByBattle)
	for result := range w.scanBoard(ctx, opt.board, opt.from, opt.to) {
		if result.Err != nil {
			return trending, result.Err
		} else {
			t.addPost(result.Post)
		}
	}

	trending.Threads = t.deviate(opt.deviate)
	return trending, nil
}
