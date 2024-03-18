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

func TestParsePath(t *testing.T) {
	ok := reflect.DeepEqual(parsePath("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePath("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePath("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePath failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/hello/geektutu")
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}
	if n.path != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if params["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.path, params["name"])
}

func TestGetRoute2(t *testing.T) {
	r := newTestRouter()
	n1, params1 := r.getRoute("GET", "/assets/file1.txt")
	ok1 := n1.path == "/assets/*filePath" && params1["filePath"] == "file1.txt"
	if !ok1 {
		t.Fatal("path should be /assets/*filePath & filePath should be file1.txt")
	}
	n2, params2 := r.getRoute("GET", "/assets/css/test.css")
	ok2 := n2.path == "/assets/*filePath" && params2["filePath"] == "css/test.css"
	if !ok2 {
		t.Fatal("path should be /assets/*filePath & filePath should be css/test.css")
	}
}

func TestGetRoutes(t *testing.T) {
	r := newTestRouter()
	nodes := r.getRoutes("GET")
	for i, n := range nodes {
		fmt.Println(i+1, n)
	}
	if len(nodes) != 5 {
		t.Fatal("the number of routes should be 5")
	}
}
