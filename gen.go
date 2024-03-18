package gen

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
	// for html render
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default use Logger & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger())
	return engine
}

// SetFuncMap assume that only invoke one time
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(path string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(path)) // batch parse, like *.tmpl
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

func (g *RouterGroup) addRoute(method, comp string, handler ...HandlerFunc) {
	path_ := g.prefix + comp
	log.Printf("Route %4s - %s\n", method, path_)
	g.engine.router.addRoute(method, path_, handler...)
}

func (g *RouterGroup) GET(path string, handler ...HandlerFunc) {
	g.addRoute("GET", path, handler...)
}

func (g *RouterGroup) POST(path string, handler ...HandlerFunc) {
	g.addRoute("POST", path, handler...)
}

func (g *RouterGroup) DELETE(path string, handler ...HandlerFunc) {
	g.addRoute("DELETE", path, handler...)
}

func (g *RouterGroup) PUT(path string, handler ...HandlerFunc) {
	g.addRoute("PUT", path, handler...)
}

func (g *RouterGroup) HEAD(path string, handler ...HandlerFunc) {
	g.addRoute("HEAD", path, handler...)
}

func (g *RouterGroup) PATCH(path string, handler ...HandlerFunc) {
	g.addRoute("PATCH", path, handler...)
}

func (g *RouterGroup) OPTIONS(path string, handler ...HandlerFunc) {
	g.addRoute("OPTIONS", path, handler...)
}

func (g *RouterGroup) CONNECT(path string, handler ...HandlerFunc) {
	g.addRoute("CONNECT", path, handler...)
}

func (g *RouterGroup) TRACE(path string, handler ...HandlerFunc) {
	g.addRoute("TRACE", path, handler...)
}

func (g *RouterGroup) Any(path string, handler ...HandlerFunc) {
	g.addRoute("GET", path, handler...)
	g.addRoute("POST", path, handler...)
	g.addRoute("DELETE", path, handler...)
	g.addRoute("PUT", path, handler...)
	g.addRoute("HEAD", path, handler...)
	g.addRoute("PATCH", path, handler...)
	g.addRoute("OPTIONS", path, handler...)
	g.addRoute("CONNECT", path, handler...)
	g.addRoute("TRACE", path, handler...)
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

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc // to find all the middlewares
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.engine = e
	c.handlers = middlewares
	e.router.handle(c)
}
