package daemon

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DaemonOption struct {
	DaemonMode bool
	PrivateKey string
	PublicKey  string
}

type Daemon struct {
	*cert
}

func NewDaemon(opt DaemonOption) *Daemon {
	return &Daemon{
		cert: newCert(opt.PublicKey, opt.PrivateKey),
	}
}

func (d *Daemon) ValidateToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 || strings.Index(auth, "Bearer ") != 0 {
		c.Abort()
		c.String(http.StatusUnauthorized, "invalid token")
		return
	}

	token := strings.Split(auth, "Bearer ")[1]
	if err := d.Validate(token); err != nil {
		c.Abort()
		c.String(http.StatusUnauthorized, "invalid token")
		return
	}

	c.Next()
}
