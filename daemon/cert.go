package daemon

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
)

type cert struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

func newCert(pub, priv string) *cert {
	pubBytes, err := ioutil.ReadFile(pub)
	fatal(err)

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	wrapFatal(err, "invalid public key")

	c := cert{
		pub: pubKey,
	}

	if priv != "" {
		privBytes, err := ioutil.ReadFile(priv)
		fatal(err)

		privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
		wrapFatal(err, "invalid private key")

		c.priv = privKey
	}

	return &c
}

func (c *cert) Sign() (string, error) {
	if c.priv == nil {
		return "", ErrValidatingCertOnly
	}

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	fatal(err)

	now := time.Now().Unix()
	t := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.StandardClaims{
		Id:        fmt.Sprintf("%x", id),
		IssuedAt:  now,
		NotBefore: now,
		Issuer:    "gossipbay-daemon",
	})

	return t.SignedString(c.priv)
}

func (c *cert) Validate(token string) error {
	if c.priv == nil {
		return ErrSigningCertOnly
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return c.pub, nil
	})
	if err != nil {
		return err
	}
	if _, ok := jwtToken.Claims.(*jwt.StandardClaims); !ok {
		return errors.Wrap(err, "invalid jwt")
	} else if !jwtToken.Valid {
		return errors.New("token expired")
	}
	return nil
}
