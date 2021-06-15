package lee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*Node
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*Node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (this *Router) handle(ctx *Context) {
	router, m := this.GetRouter(ctx.Method, ctx.Path)
	if router == nil {
		ctx.String(http.StatusNotFound, "404 NOT FOUND:%s\n", ctx.Path)
		return
	}
	if router != nil {
		ctx.Parse = m
		key := ctx.Method + "-" + router.Pattern
		ctx.handlers = append(ctx.handlers, this.handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND:%s\n", ctx.Path)
		})

	}
	ctx.Next()
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//if item[0] == '*' {
			//	break
			//}
		}
	}
	return parts
}

func (this *Router) AddRoute(method string, pattern string, handle HandlerFunc) {
	if _, ok := this.roots[method]; !ok {
		this.roots[method] = &Node{}
	}
	//插入新的节点
	this.roots[method].Insert(pattern, parsePattern(pattern), 0)
	this.handlers[method+"-"+pattern] = handle
}

func (this *Router) GetRouter(method string, pattern string) (*Node, map[string]string) {
	parts := parsePattern(pattern)
	m := make(map[string]string)
	node := this.roots[method].Search(parts, 0)
	if node != nil {
		paths := parsePattern(node.Pattern)
		for index, v := range paths {
			if v[0] == '*' || v[0] == ':' {
				m[v] = parts[index]
			}
		}
	}
	return node, m
}
