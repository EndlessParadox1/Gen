package gen

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]any

type Context struct {
	Request    *http.Request
	Writer     http.ResponseWriter
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Request: req,
		Writer:  w,
		Path:    req.URL.Path,
		Method:  req.Method,
	}
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, a ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, a)))
}

func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//func (c *Context) Deadline() (deadline time.Time, ok bool) {
//	return
//}
//
//func (c *Context) Done() <-chan struct{} {
//	return nil
//}
//
//func (c *Context) Err() error {
//	return nil
//}
//
//func (c *Context) Value(key any) any {
//	if key == 0 {
//		return c.Request
//	}
//	return nil
//}
