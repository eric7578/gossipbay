package schedule

import (
	"fmt"
	"time"

	"github.com/eric7578/gossipbay/repo"
)

func Prune(r repo.Repository, user string, from time.Time, to time.Time, labels ...string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("prune failed %s", e)
		}
	}()

	comments := r.ListComments(from)
	for _, c := range comments {
		if c.Author == user && (to.IsZero() || c.UpdatedAt.Before(to)) {
			r.RemoveComment(c.ID)
		}
	}

	return nil
}
