package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseArgs(t *testing.T) {
	args := parseArgs("pttweb board=Gossiping deviate=0.8")

	assert.Equal(t, args["_type"], "pttweb")
	assert.Equal(t, args["board"], "Gossiping")
	assert.Equal(t, args["deviate"], "0.8")
}
