package myhttp

import (
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
	requestLine RequestLine
	headers     map[string]string
	body        []byte
}

type RequestLine struct {
	httpMethod HttpMethod
	target     string
	version    string
}

func FromReader(r io.Reader) (*HTTPRequest, error) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	_, err := r.Read(buffer.Bytes())
	if err != nil {
		err := fmt.Errorf("Cannot read bytes from the reader. Failed with error: %v", err)
		return nil, err
	}

	dataStr := buffer.String()
	//to split between the payload and the header part
	sections := strings.Split(dataStr, "\r\n\r\n")
	headerPart := sections[0]
	payloadPart := sections[1]

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
		parts := strings.SplitN(line, ":", 1)
		headers[parts[0]] = parts[1]
	}

	return &HTTPRequest{
		requestLine: RequestLine{
			httpMethod: method,
			target:     requestLineStr[1],
			version:    requestLineStr[2],
		},
		headers: headers,
		body:    []byte(payloadPart),
	}, nil
}

func (hr *HTTPRequest) String() string {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	var method string
	switch hr.requestLine.httpMethod {
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
	buffer.WriteString(hr.requestLine.target + "\n")
	buffer.WriteString(hr.requestLine.version + "\n")

	buffer.WriteString("\n")

	for key, value := range hr.headers {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	buffer.WriteString("\n")

	buffer.Write(hr.body)

	return buffer.String()
}
