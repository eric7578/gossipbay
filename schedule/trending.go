package schedule

import (
	"context"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/eric7578/gossipbay/crawler"
)

type TrendingOption struct {
	Board   string
	From    time.Time
	To      time.Time
	Timeout time.Duration
	Deviate float64
}

func (opt TrendingOption) isValid() bool {
	return opt.Board != "" && !opt.From.IsZero()
}

type Thread struct {
	Score   float64        `json:"score"`
	Deviate float64        `json:"deviate"`
	Posts   []crawler.Post `json:"posts"`
}

func (t *Thread) sortPosts() {

}

func ScoreByBattle(p crawler.Post) float64 {
	return float64(p.BattlePush)
}

type trending struct {
	scoreFunc func(p crawler.Post) float64
	threads   map[string]*Thread
}

func newTrending(scoreFunc func(p crawler.Post) float64) *trending {
	return &trending{
		scoreFunc: scoreFunc,
		threads:   make(map[string]*Thread),
	}
}

func (t *trending) addPost(p crawler.Post) {
	group := genGroup(p.Title)
	thread, ok := t.threads[group]
	if ok {
		thread.Score += t.scoreFunc(p)
		thread.Posts = append(thread.Posts, p)
	} else {
		thread = &Thread{
			Score: t.scoreFunc(p),
			Posts: []crawler.Post{p},
		}
		t.threads[group] = thread
	}
}

func (t *trending) deviate(threshold float64) []*Thread {
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

	threads := make([]*Thread, 0)
	offset := maxScore - minScore
	for _, thread := range t.threads {
		thread.Deviate = thread.Score / offset
		if threshold <= thread.Deviate {
			thread.sortPosts()
			threads = append(threads, thread)
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

func (s *Scheduler) Trending(opt TrendingOption) ([]*Thread, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if opt.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), opt.Timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	t := newTrending(ScoreByBattle)
	for result := range s.crawler.ScanBoard(ctx, opt.Board, opt.From, opt.To) {
		if result.Err != nil {
			return nil, result.Err
		} else {
			t.addPost(result.Post)
		}
	}
	return t.deviate(opt.Deviate), nil
}
