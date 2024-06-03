package myhttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type HttpMethod int

const (
	GET HttpMethod = iota
	POST
	PUT
	PATCH
	DELETE
	UNKNOWN
)

type HTTPRequest struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpMethod HttpMethod
	Target     string
	Version    string
}

func FromReader(r io.Reader) (*HTTPRequest, error) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	reader := bufio.NewReader(r)
	tmp := make([]byte, 2048)
	n, err := reader.Read(tmp)
	if err != nil {
		err := fmt.Errorf("There was some error while reading from the socket. err: %v\n", err)
		return nil, err
	}
	buffer.Write(tmp[:n])

	dataStr := buffer.String()
	sections := strings.Split(dataStr, "\r\n\r\n")
	headerPart := sections[0]

	var payloadPart string
	if len(sections) > 1 {
		payloadPart = sections[1]
	}

	headerSections := strings.Split(headerPart, "\r\n")
	if len(sections) == 0 {
		err := fmt.Errorf("header sections is empty.")
		return nil, err
	}

	requestLineStr := strings.Split(headerSections[0], " ")

	var method HttpMethod
	switch requestLineStr[0] {
	case "GET":
		method = GET
	case "PUT":
		method = PUT
	case "PATCH":
		method = PATCH
	case "POST":
		method = POST
	case "DELETE":
		method = DELETE
	default:
		method = UNKNOWN
	}

	headers := make(map[string]string)
	for _, line := range headerSections[1:] {
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	return &HTTPRequest{
		RequestLine: RequestLine{
			HttpMethod: method,
			Target:     requestLineStr[1],
			Version:    requestLineStr[2],
		},
		Headers: headers,
		Body:    []byte(payloadPart),
	}, nil
}

func (hr *HTTPRequest) String() string {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	var method string
	switch hr.RequestLine.HttpMethod {
	case GET:
		method = "GET"
	case PUT:
		method = "PUT"
	case POST:
		method = "POST"
	case PATCH:
		method = "PATCH"
	case DELETE:
		method = "DELETE"
	case UNKNOWN:
		method = "UNKNOWN"
	}

	buffer.WriteString(method + "\n")
	buffer.WriteString(hr.RequestLine.Target + "\n")
	buffer.WriteString(hr.RequestLine.Version + "\n")

	buffer.WriteString("\n")

	for key, value := range hr.Headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	buffer.WriteString("\n")

	buffer.Write(hr.Body)

	return buffer.String()
}

func (hr *HTTPRequest) PathParams() []string {
	params := strings.Split(hr.RequestLine.Target, "/")
	if len(params) > 0 && params[len(params)-1] == "" {
		params = params[:len(params)-1]
	}
	return params
}
