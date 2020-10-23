package utils

import "time"

var (
	Taipei = mustLoadLocation()
)

func mustLoadLocation() *time.Location {
	taipei, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	return taipei
}
