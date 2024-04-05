package api

import "github.com/emicklei/go-restful/v3"

// 写一个 go restful 框架的Handler
// 使用的v3版本
//
//	Gin （ctx *gin.Context)
type Info struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

// GoRestful (r *restful.Request, w *restful.Response)
func (h *handler) CreateSecret(r *restful.Request, w *restful.Response) {
	// 1. 获取请求 Request Entity

	// 2. Gin Bind
	// r.ReadEntity()

	// 2. 处理请求

	// 3. 返回请求
	// c.JSON()
	// w.WriteError()
	outPut := &Info{
		Data: map[string]any{"hh": 0},
		Code: 0,
		Msg:  "success",
	}
	w.WriteAsJson(outPut)
}
