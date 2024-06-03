package myhttp

import "testing"

func TestRouterTreeGood(t *testing.T) {
	rt := NewRouterTree()
	rt.Add("/", func(hc *HttpContext) HttpResponse {
		return *NewHttpResponseBuilder().
			WithBody([]byte("Root")).
			Build()
	})

	rt.Add("/hello", func(hc *HttpContext) HttpResponse {
		return *NewHttpResponseBuilder().
			WithBody([]byte("Hola")).
			Build()
	})

	rt.Add("/echo/:data", func(hc *HttpContext) HttpResponse {
		return *NewHttpResponseBuilder().
			WithBody([]byte(hc.pathParams["data"])).
			Build()
	})

	rt.Add("/:first/:second/:third", func(hc *HttpContext) HttpResponse {
		return *NewHttpResponseBuilder().
			WithBody([]byte("LotsOfPathParams")).
			Build()
	})

	c1 := NewHttpContext(&HttpRequest{
		RequestLine: RequestLine{
			Target: "/",
		},
	})
	f1, err := rt.GetHandler(c1)

	if err != nil {
		t.Fatalf("Root Handler could not be fetched: %v", err)
	}

	resp1 := f1(c1)
	if string(resp1.Body) != "Root" {
		t.Fatalf("Root Body not equal. Want: Root, Got: %s", string(resp1.Body))
	}

	c2 := NewHttpContext(&HttpRequest{
		RequestLine: RequestLine{
			Target: "/hello",
		},
	})
	f2, err := rt.GetHandler(c2)

	if err != nil {
		t.Fatalf("Hello Handler could not be fetched: %v", err)
	}

	resp2 := f2(c2)
	if string(resp2.Body) != "Hola" {
		t.Fatalf("Hello Body not equal. Want: Root, Got: %s", string(resp1.Body))
	}

	c3 := NewHttpContext(&HttpRequest{
		RequestLine: RequestLine{
			Target: "/echo/omegalul",
		},
	})
	f3, err := rt.GetHandler(c3)

	if err != nil {
		t.Fatalf("Hello Handler could not be fetched: %v", err)
	}

	resp3 := f3(c3)
	if string(resp3.Body) != "omegalul" {
		t.Fatalf("Echo Body not equal. Want: omegalul, Got: %s", string(resp1.Body))
	}

	c4 := NewHttpContext(&HttpRequest{
		RequestLine: RequestLine{
			Target: "/this/is/cool",
		},
	})
	f4, err := rt.GetHandler(c4)

	if err != nil {
		t.Fatalf("Hello Handler could not be fetched: %v", err)
	}

	resp4 := f4(c4)
	if string(resp4.Body) != "LotsOfPathParams" {
		t.Fatalf("Param Test Body not equal. Want: LotsOfPathParams, Got: %s", string(resp1.Body))
	}

	if c4.pathParams["first"] != "this" {
		t.Fatalf("Path Params are wrong. Want: this, Got: %s", c4.pathParams["first"])
	}
	if c4.pathParams["second"] != "is" {
		t.Fatalf("Path Params are wrong. Want: is, Got: %s", c4.pathParams["first"])
	}
	if c4.pathParams["third"] != "cool" {
		t.Fatalf("Path Params are wrong. Want: cool, Got: %s", c4.pathParams["first"])
	}
}

func TestRouterTreeBad(t *testing.T) {
	rt := NewRouterTree()

	c := NewHttpContext(&HttpRequest{
		RequestLine: RequestLine{
			Target: "Omega",
		},
	})
	_, err := rt.GetHandler(c)
	if err == nil {
		t.Fatalf("Error is not nil when the route is absent.")
	}
}
