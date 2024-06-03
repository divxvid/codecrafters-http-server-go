package myhttp

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
