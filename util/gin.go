package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetParams(c *gin.Context) (content map[string]interface{}) {
	switch {
	case c.Request.Method == "GET":
		content = GetQueryParams(c)
	case c.ContentType() == "application/json":
		content = GetBodyParams(c)
	default:
		content = GetPostFormParams(c)
	}
	return
}

func GetQueryParams(c *gin.Context) (params map[string]interface{}) {
	if params, ok := c.Get("queryParams"); ok {
		return params.(map[string]interface{})
	}

	params = make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		params[k] = v[0]
	}

	c.Set("queryParams", params)
	return
}

func GetBodyParams(c *gin.Context) (params map[string]interface{}) {
	if params, ok := c.Get("bodyParams"); ok {
		return params.(map[string]interface{})
	}

	params = make(map[string]interface{})
	data, _ := c.GetRawData()

	//把读过的字节流重新放到body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	json.Unmarshal(data, &params)

	c.Set("bodyParams", params)
	return
}

func GetPostFormParams(c *gin.Context) (params map[string]interface{}) {
	if params, ok := c.Get("postFormParams"); ok {
		return params.(map[string]interface{})
	}

	params = make(map[string]interface{})

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return
		}
	}

	for k, v := range c.Request.PostForm {
		if len(v) > 1 {
			params[k] = v
		} else if len(v) == 1 {
			params[k] = v[0]
		}
	}

	c.Set("postFormParams", params)
	return
}
