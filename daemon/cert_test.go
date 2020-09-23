package daemon

import (
	"path"
	"testing"

	"github.com/eric7578/gossipbay/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_cert_Sign_Validate(t *testing.T) {
	pub := path.Join(testutil.MustGetwd(), "testdata/id_rsa.pub")
	priv := path.Join(testutil.MustGetwd(), "testdata/id_rsa")
	c := newCert(pub, priv)

	token, signErr := c.Sign()
	validateErr := c.Validate(token)
	invalidateErr := c.Validate("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI0NzQ3YjNkMDgwMDAxNjgiLCJpYXQiOjE2MDA4NzA2MjcsImlzcyI6Imdvc3NpcGJheS1kYWVtb24iLCJuYmYiOjE2MDA4NzA2Mjd9.Bza0X4v476rLPwNTla9oHJg63nVFnodZm0vxgOr00Ji3Hi5dtu0r1-5e_mdweMJu27INih5jZblRzK2rYC1ymA")

	assert.NotEqual(t, "", token)
	assert.Nil(t, signErr)
	assert.Nil(t, validateErr)
	assert.NotNil(t, invalidateErr)
}
