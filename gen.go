package gen

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup // for group nesting
	engine      *Engine      // all groups share the same engine, to access its `router`
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // stores all groups
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
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

//func (g *RouterGroup) Use(middleware HandlerFunc) {
//
//}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	path := g.prefix + comp
	log.Printf("Route %4s - %s\n", method, path)
	g.engine.router.addRoute(method, path, handler)
}

func (g *RouterGroup) GET(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler)
}

func (g *RouterGroup) POST(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
