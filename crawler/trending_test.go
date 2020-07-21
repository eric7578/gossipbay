package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrending_addPosts(t *testing.T) {
	tr := NewTrending()
	tr.addPosts(
		Post{Title: "Title", NumPush: 3},
		Post{Title: "Re: Title", NumPush: 5},
		Post{Title: "Title2", NumPush: 1},
		Post{Title: "Title3", NumPush: 7},
	)
	threads := tr.Deviate(0.8)

	assert.Equal(t, 2, len(threads))
	assert.Equal(t, "Title", threads[0].Title)
	assert.Equal(t, "Title3", threads[1].Title)
}
