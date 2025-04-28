package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string //前缀
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + prefix,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute是向特定的路由组添加新的路由规则
// 将group.prefix与comp组合成为一个完整的pattern
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 添加一个处理pattern路径的GET请求的路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 添加一个处理pattern路径的POST请求的路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

//请求处理的过程：
//1.客户端向服务器发送http请求
//2.服务器端接收请求，http.Server调用engine的ServerHTTP函数
//3.ServeHTTP调用newContext封装请求和响应
//4.ServeHTTPP将Context传给router的handle方法
//5.router的handle方法根据请求的路径和方法，在路由树中查找匹配的节点，并调用相应的处理函数
//6.处理函数通过Context生成响应返回给客户端
