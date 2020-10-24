package daemon

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Daemon) createJob(c *gin.Context) {
	body := struct {
		Type string
		Args map[string]string
	}{}
	c.BindJSON(&body)

	res, err := d.crawler.CreateJob(body.Type, body.Args)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}
