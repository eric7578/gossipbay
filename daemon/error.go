package daemon

import "github.com/pkg/errors"

var (
	ErrSigningCertOnly    = errors.New("cert can only be used in signing")
	ErrValidatingCertOnly = errors.New("cert can only be used in validating")
)

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func wrapFatal(err error, message string) {
	if err != nil {
		panic(errors.Wrap(err, message))
	}
}
