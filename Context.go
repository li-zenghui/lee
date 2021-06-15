package lee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	Parse      map[string]string
	//宽展中间件入口
	handlers []HandlerFunc
	index    int //执行到第几个
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

func (this *Context) Next() {
	this.index++
	s := len(this.handlers)
	for ; this.index < s; this.index++ {
		this.handlers[this.index](this)
	}

}

func (this *Context) PostFrom(key string) string {
	value := this.Req.FormValue(key)
	return value
}

func (this *Context) Query(key string) string {
	return this.Req.URL.Query().Get(key)
}

func (this *Context) SetStatus(code int) {
	this.StatusCode = code
	//this.Writer.WriteHeader(code)
}

func (this *Context) SetHeader(key string, value string) {
	this.Writer.Header().Set(key, value)
}

func (this *Context) String(code int, format string, values ...interface{}) {
	this.SetHeader("Content-Type", "text/plain")
	this.SetStatus(code)
	this.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (this *Context) JSON(code int, obj interface{}) {
	this.SetHeader("Content-Type", "application/json")
	this.SetStatus(code)
	encoder := json.NewEncoder(this.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(this.Writer, err.Error(), 500)
	}
}

func (this *Context) Data(code int, data []byte) {
	this.SetStatus(code)
	this.Writer.Write(data)

}

func (this *Context) HTML(code int, html string) {
	this.SetHeader("Content-Type", "text/html")
	this.SetStatus(code)
	this.Writer.Write([]byte(html))
}
