package myhttp

import "strings"

type HttpContext struct {
	request    *HttpRequest
	pathParams map[string]string
}

func NewHttpContext(req *HttpRequest) *HttpContext {
	return &HttpContext{
		request:    req,
		pathParams: make(map[string]string),
	}
}

func (ctx *HttpContext) PathParam(key string) string {
	return ctx.pathParams[key]
}

func (ctx *HttpContext) GetRequestHeader(key string) string {
	//lowercase the key
	key = strings.ToLower(key)
	return ctx.request.Headers[key]
}

func (ctx *HttpContext) GetRequestBody() []byte {
	return ctx.request.Body
}
