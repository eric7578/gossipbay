package ptt

import (
	"context"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

type TrendingOption struct {
	Board   string        `json:"board"`
	From    time.Time     `json:"trendingFrom"`
	To      time.Time     `json:"trendingTo"`
	Timeout time.Duration `json:"timeout"`
	Deviate float64       `json:"deviate"`
}

func (opt TrendingOption) IsValid() bool {
	return opt.Board != "" && !opt.From.IsZero()
}

type Trending struct {
	TrendingOption
	StartAt time.Time `json:"jobStartAt"`
	EndAt   time.Time `json:"jobEndAt"`
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

func (c *PageCrawler) trending(ctx context.Context, opt TrendingOption) (Trending, error) {
	var (
		cancel   context.CancelFunc
		trending = Trending{
			TrendingOption: opt,
			StartAt:        time.Now(),
		}
	)
	if opt.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opt.Timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	t := newTrending(ScoreByBattle)
	for result := range c.ScanBoard(ctx, opt.Board, opt.From, opt.To) {
		if result.Err != nil {
			return trending, result.Err
		} else {
			t.addPost(result.Post)
		}
	}

	trending.Threads = t.deviate(opt.Deviate)
	trending.EndAt = time.Now()
	return trending, nil
}

func (s *PageCrawler) Trending(ctx context.Context, opts ...TrendingOption) ([]Trending, error) {
	fad := make([]Trending, 0)
	trendingc := make(chan Trending)
	done := make(chan struct{})

	go func() {
		for t := range trendingc {
			fad = append(fad, t)
		}
		done <- struct{}{}
	}()

	var wg sync.WaitGroup
	wg.Add(len(opts))
	for _, opt := range opts {
		go func(opt TrendingOption) {
			defer wg.Done()
			t, err := s.trending(ctx, opt)
			if err != nil {
				return
			}
			trendingc <- t
		}(opt)
	}
	wg.Wait()
	close(trendingc)
	<-done
	return fad, nil
}
