package gen

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

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
	log.SetPrefix("[GEN] ")
	return engine
}

// Default use Logger & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// SetFuncMap assume that only invoke one time
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(path string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(path)) // like *.tmpl
}

func (e *Engine) LoadHTMLFiles(files ...string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseFiles(files...)) // like a.tmpl, b.tmpl...
}

func (e *Engine) Run(addr string) error {
	log.Printf("Listening and serving HTTP on %s\n", addr)
	return http.ListenAndServe(addr, e)
}

func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	log.Printf("Listening and serving HTTPS on %s\n", addr)
	return http.ListenAndServeTLS(addr, certFile, keyFile, e)
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
