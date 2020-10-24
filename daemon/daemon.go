package daemon

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DaemonOption struct {
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

func (d *Daemon) validateToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 || strings.Index(auth, "Bearer ") != 0 {
		c.Abort()
		c.String(http.StatusUnauthorized, "token required")
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

func (d *Daemon) Run(port string) {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.Use(d.validateToken)
		api.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		api.POST("/job", d.createJob)
	}

	r.Run(port)
}
