package gen

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/quic-go/quic-go/http3"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router *httprouter.Router
	groups []*RouterGroup // stores all groups
	// for html render
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: httprouter.New()}
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

func (e *Engine) RunQUIC(addr, certFile, keyFile string) error {
	log.Printf("Listening and serving HTTP3 on %s\n", addr)
	return http3.ListenAndServeQUIC(addr, certFile, keyFile, e)
}

func (e *Engine) RunHttp3(addr, certFile, keyFile string) error {
	log.Printf("Listening and serving HTTP3 on %s\n", addr)
	return http3.ListenAndServeTLS(addr, certFile, keyFile, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.ServeHTTP(w, req)
}
