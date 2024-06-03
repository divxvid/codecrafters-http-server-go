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
