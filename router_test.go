package gen

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filePath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePath("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePath("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePath("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geektutu")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.path != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.path, ps["name"])
}

func TestGetRoute2(t *testing.T) {
	r := newTestRouter()
	n1, ps1 := r.getRoute("GET", "/assets/file1.txt")
	ok1 := n1.path == "/assets/*filePath" && ps1["filePath"] == "file1.txt"
	if !ok1 {
		t.Fatal("path should be /assets/*filePath & filePath should be file1.txt")
	}
	n2, ps2 := r.getRoute("GET", "/assets/css/test.css")
	ok2 := n2.path == "/assets/*filePath" && ps2["filePath"] == "css/test.css"
	if !ok2 {
		t.Fatal("path should be /assets/*filePath & filePath should be css/test.css")
	}
}
