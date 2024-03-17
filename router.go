package gen

import "strings"

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// only one * is allowed
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

func (r *router) addRoute(method string, path string, handler HandlerFunc) {
	parts := parsePath(path)
	key := method + "-" + path
	root, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	root.insert(path, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	var params map[string]string
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
		key := c.Method + "-" + c.Path
		r.handlers[key](c)
	} else {
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
}
