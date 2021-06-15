package lee

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	*RouterGroup
	router *Router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = append(engine.groups, engine.RouterGroup)
	return engine
}

//将路由加入映射表
func (this *Engine) AddRouter(method string, pattern string, handler HandlerFunc) {
	this.router.AddRoute(method, pattern, handler)
}

func (this *Engine) Get(pattern string, handler HandlerFunc) {
	this.AddRouter("GET", pattern, handler)
}

func (this *Engine) Post(pattern string, handler HandlerFunc) {
	this.AddRouter("POST", pattern, handler)
}

func (this *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, this)
}

func (this *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var newMiddleWare []HandlerFunc
	c := NewContext(w, req)
	for _, group := range this.groups {
		if strings.HasPrefix(c.Path, group.Prefix) {
			newMiddleWare = append(newMiddleWare, group.Middleware...)
		}
	}
	c.handlers = newMiddleWare
	this.router.handle(c)
}



type RouterGroup struct {
	Prefix     string        //前缀
	Middleware []HandlerFunc //中间件插口
	Parent     *RouterGroup  //父分组
	engine     *Engine       //引擎
}

func (this *RouterGroup) Group(prefix string) *RouterGroup {
	//添加分组
	engine := this.engine
	newGroup := &RouterGroup{
		engine: engine,
		Prefix: this.Prefix + prefix,
		Parent: this,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//封装路由方法

func (this *RouterGroup) AddRouter(method string, comp string, handle HandlerFunc) {
	pattern := this.Prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	this.engine.router.AddRoute(method, pattern, handle)
}

func (this *RouterGroup) Get(pattern string, handle HandlerFunc) {
	this.AddRouter("GET", pattern, handle)
}

func (this *RouterGroup) Post(pattern string, handle HandlerFunc) {
	this.AddRouter("POST", pattern, handle)
}

func (this *RouterGroup)Use(middleware... HandlerFunc)  {
	this.Middleware = append(this.Middleware, middleware...)
}



func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				//c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

func onlyForV2() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
		c.Next()
		log.Printf("访问结束！")
	}
}

func Default() *Engine {
	engine:= New()
	engine.Use(onlyForV2(),Recovery())
	return engine
}