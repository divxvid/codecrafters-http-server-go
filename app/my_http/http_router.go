package myhttp

type HandlerFunc func(*HttpContext) HttpResponse

type Router struct {
	routes map[string]RouterTree
}

func NewRouter() *Router {
	routes := make(map[string]RouterTree)
	routes["GET"] = NewRouterTree()
	routes["POST"] = NewRouterTree()
	routes["PATCH"] = NewRouterTree()
	routes["PUT"] = NewRouterTree()
	routes["DELETE"] = NewRouterTree()
	return &Router{
		routes: routes,
	}
}

func (r *Router) GET(path string, f HandlerFunc) error {
	return r.routes["GET"].Add(path, f)
}

func (r *Router) POST(path string, f HandlerFunc) error {
	return r.routes["POST"].Add(path, f)
}

func (r *Router) PATCH(path string, f HandlerFunc) error {
	return r.routes["PATCH"].Add(path, f)
}

func (r *Router) PUT(path string, f HandlerFunc) error {
	return r.routes["PUT"].Add(path, f)
}

func (r *Router) DELETE(path string, f HandlerFunc) error {
	return r.routes["DELETE"].Add(path, f)
}

func (r *Router) HandleRequest(request *HttpRequest) (*HttpResponse, error) {
	ctx := NewHttpContext(request)
	f, err := r.routes[request.RequestLine.HttpMethod].GetHandler(ctx)
	if err != nil {
		return nil, err
	}
	response := f(ctx)
	return &response, nil
}
