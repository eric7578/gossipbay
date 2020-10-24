package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Daemon) createJob(c *gin.Context) {
	body := struct {
		Type     string
		Args     map[string]string
		Callback string
	}{}
	c.BindJSON(&body)

	if body.Callback == "" {
		res, err := d.crawler.CreateJob(body.Type, body.Args)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, res)
		return
	}

	go func() {
		res, _ := d.crawler.CreateJob(body.Type, body.Args)
		d.doCallback(body.Callback, res)
	}()
	c.Status(http.StatusOK)
}

func (d *Daemon) doCallback(url string, body interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	c := http.Client{}
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("request to %s got status code error: [%d] %s", url, res.StatusCode, res.Status)
	}
	return nil
}
