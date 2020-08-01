package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrending_addPosts(t *testing.T) {
	posts := []Post{
		{Title: "Title", NumPush: 3},
		{Title: "Re: Title", NumPush: 5},
		{Title: "Title2", NumPush: 1},
		{Title: "Title3", NumPush: 7},
	}
	tr := NewTrending(posts)
	threads := tr.Deviate(0.8)

	assert.Equal(t, 2, len(threads))
	assert.Equal(t, "Title", threads[0].Title)
	assert.Equal(t, "Title3", threads[1].Title)
	assert.Equal(t, 8, threads[0].NumPush)
	assert.Equal(t, 7, threads[1].NumPush)
}
