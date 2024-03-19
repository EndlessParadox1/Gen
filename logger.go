package gen

import (
	"log"
	"time"
)

// Logger middleware that top group use
func Logger() HandlerFunc {
	return func(c *Context) {
		log.SetPrefix("[GEN] ")
		t := time.Now()
		c.Next()
		log.Printf("| %d | %13v | %15s | %-7s  %s\n", c.StatusCode, time.Since(t), c.RemoteIP(), c.Method, c.Path)
	}
}
