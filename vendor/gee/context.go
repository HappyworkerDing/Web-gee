package gee

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
	Params     map[string]string
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// format是一个站位符，提供一个格式化字符串的模版
// eg.c.String(200, "Hello, %s!", "Alice") Alice就是values 输出的结果是Helli Alice ， 其中format就是%s
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// obj代表的是需要序列化为JSON格式的数据对象
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	//josn包的NewEncoder函数创建一个JSON编码器 *Encoder json.Encoder的主要功能是将go数据结构编码成JSON格式
	//func NewEncoder(w io.Writer) *Encoder
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		//http.Error: func Error(w ResponseWriter , error string , code int)
		//将错误信息写入HTTP响应体，并设置错误状态码
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
