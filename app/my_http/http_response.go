package myhttp

import (
	"bytes"
	"fmt"
	"net"
)

type HttpResponse struct {
	StatusLine StatusLine
	Headers    map[string]string
	Body       []byte
}

type StatusLine struct {
	Protocol   string
	StatusCode int
	StatusText string
}

type HttpResponseBuilder struct {
	Protocol   string
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func NewHttpResponseBuilder() *HttpResponseBuilder {
	return &HttpResponseBuilder{
		Protocol:   "HTTP/1.1",
		StatusCode: 200,
		StatusText: "OK",
		Headers:    make(map[string]string),
		Body:       nil,
	}
}

func (b *HttpResponseBuilder) WithProtocol(protocol string) *HttpResponseBuilder {
	b.Protocol = protocol
	return b
}

func (b *HttpResponseBuilder) WithStatusCode(code int) *HttpResponseBuilder {
	b.StatusCode = code
	return b
}

func (b *HttpResponseBuilder) WithStatusText(text string) *HttpResponseBuilder {
	b.StatusText = text
	return b
}

func (b *HttpResponseBuilder) WithHeader(key, value string) *HttpResponseBuilder {
	b.Headers[key] = value
	return b
}

func (b *HttpResponseBuilder) WithBody(body []byte) *HttpResponseBuilder {
	b.Body = body
	return b
}

func (b *HttpResponseBuilder) Build() *HttpResponse {
	return &HttpResponse{
		StatusLine: StatusLine{
			Protocol:   b.Protocol,
			StatusCode: b.StatusCode,
			StatusText: b.StatusText,
		},
		Headers: b.Headers,
		Body:    b.Body,
	}
}

func (hr *HttpResponse) WriteToConn(conn net.Conn) (n int, err error) {
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
