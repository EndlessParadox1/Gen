package gen

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v\n", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
