package myhttp

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	StatusLine StatusLine
	Headers    map[string]string
	Body       []byte
}

type StatusLine struct {
	Protocol   string
	StatusCode int
	StatusText string
}

type ResponseBuilder struct {
	protocol   string
	statusCode int
	statusText string
	headers    map[string]string
	body       []byte
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		protocol:   "HTTP/1.1",
		statusCode: 200,
		statusText: "OK",
		headers:    make(map[string]string),
		body:       nil,
	}
}

func (b *ResponseBuilder) Protocol(protocol string) *ResponseBuilder {
	b.protocol = protocol
	return b
}

func (b *ResponseBuilder) StatusCode(code int) *ResponseBuilder {
	b.statusCode = code
	b.statusText = http.StatusText(code)
	return b
}

func (b *ResponseBuilder) StatusText(text string) *ResponseBuilder {
	b.statusText = text
	return b
}

func (b *ResponseBuilder) Header(key, value string) *ResponseBuilder {
	b.headers[key] = value
	return b
}

func (b *ResponseBuilder) Body(body []byte) *ResponseBuilder {
	b = b.Header("Content-Length", fmt.Sprintf("%d", len(body)))
	b.body = body
	return b
}

func (b *ResponseBuilder) Build() *Response {
	return &Response{
		StatusLine: StatusLine{
			Protocol:   b.protocol,
			StatusCode: b.statusCode,
			StatusText: b.statusText,
		},
		Headers: b.headers,
		Body:    b.body,
	}
}

func (hr *Response) WriteToConn(conn net.Conn) (n int, err error) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	buffer.WriteString(hr.StatusLine.Protocol + " ")
	buffer.WriteString(fmt.Sprintf("%d", hr.StatusLine.StatusCode) + " ")
	buffer.WriteString(hr.StatusLine.StatusText)
	buffer.WriteString("\r\n")

	for key, value := range hr.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	buffer.WriteString("\r\n")

	if hr.Body != nil {
		buffer.Write(hr.Body)
	}

	n, err = conn.Write(buffer.Bytes())
	return
}
