package gen

import (
	"log"
	"net/http"
	"path"
	"reflect"
	"runtime"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup // for group nesting
	engine      *Engine      // all groups share the same engine, to access its `router`
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) addRoute(method, comp string, handlers ...HandlerFunc) {
	path_ := g.prefix + comp
	len_ := len(handlers)
	f := handlers[len_-1]
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	log.Printf("%-6s %-25s --> %s (%d handlers)\n", method, path_, name, len_)
	g.engine.router.addRoute(method, path_, handlers...)
}

func (g *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodGet, path, handlers...)
}

func (g *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodPost, path, handlers...)
}

func (g *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodDelete, path, handlers...)
}

func (g *RouterGroup) HEAD(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodHead, path, handlers...)
}

func (g *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodPatch, path, handlers...)
}

func (g *RouterGroup) CONNECT(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodConnect, path, handlers...)
}

func (g *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodOptions, path, handlers...)
}

func (g *RouterGroup) TRACE(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodTrace, path, handlers...)
}

func (g *RouterGroup) Any(path string, handlers ...HandlerFunc) {
	g.addRoute(http.MethodGet, path, handlers...)
	g.addRoute(http.MethodPost, path, handlers...)
	g.addRoute(http.MethodDelete, path, handlers...)
	g.addRoute(http.MethodPut, path, handlers...)
	g.addRoute(http.MethodHead, path, handlers...)
	g.addRoute(http.MethodPatch, path, handlers...)
	g.addRoute(http.MethodConnect, path, handlers...)
	g.addRoute(http.MethodOptions, path, handlers...)
	g.addRoute(http.MethodTrace, path, handlers...)
}

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		filePath := c.Param("filePath")
		if _, err := fs.Open(filePath); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (g *RouterGroup) Static(relativePath, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filePath")
	g.GET(urlPath, handler)
}
