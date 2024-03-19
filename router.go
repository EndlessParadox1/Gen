package gen

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string][]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string][]HandlerFunc),
	}
}

// Only one * is allowed
func parsePath(path string) (parts []string) {
	items := strings.Split(path, "/")
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return
}

func (r *router) addRoute(method, path string, handlers ...HandlerFunc) {
	parts := parsePath(path)
	key := method + "-" + path
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handlers
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	params := make(map[string]string)
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePath(n.path)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.path
		c.handlers = append(c.handlers, r.handlers[key]...)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
