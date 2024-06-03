package myhttp

import "fmt"

type HandlerFunc func(*HttpContext) HttpResponse

type Router struct {
	routes map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) handlePathParams(path string) {

}

func (r *Router) GET(path string, f HandlerFunc) {
	path = fmt.Sprintf("GET %s", path)
	r.routes[path] = f
}

func (r *Router) POST(path string, f HandlerFunc) {
	path = fmt.Sprintf("POST %s", path)
	r.routes[path] = f
}

func (r *Router) PATCH(path string, f HandlerFunc) {
	path = fmt.Sprintf("PATCH %s", path)
	r.routes[path] = f
}

func (r *Router) PUT(path string, f HandlerFunc) {
	path = fmt.Sprintf("PUT %s", path)
	r.routes[path] = f
}

func (r *Router) DELETE(path string, f HandlerFunc) {
	path = fmt.Sprintf("DELETE %s", path)
	r.routes[path] = f
}

func (r *Router) HandleRequest(request *HttpRequest) HttpResponse {
	routeKey := fmt.Sprintf(
		"%s %s",
		request.RequestLine.HttpMethod,
		request.RequestLine.Target,
	)

	f := r.routes[routeKey]
	ctx := NewHttpContext(request)
	response := f(ctx)
	return response
}
