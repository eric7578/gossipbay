package schedule

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

func Pipe(src interface{}, dest io.Writer) error {
	if bytes, err := json.Marshal(src); err != nil {
		return errors.Wrap(err, "cannot format as json")
	} else if _, err := dest.Write(bytes); err != nil {
		return errors.Wrap(err, "output failed")
	}
	return nil
}
