package myhttp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type HttpRequest struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpMethod string
	Target     string
	Version    string
}

func FromReader(r io.Reader) (*HttpRequest, error) {
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

	headers := make(map[string]string)
	for _, line := range headerSections[1:] {
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		//lowercase the key
		key = strings.ToLower(key)
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	return &HttpRequest{
		RequestLine: RequestLine{
			HttpMethod: strings.TrimSpace(requestLineStr[0]),
			Target:     strings.TrimSpace(requestLineStr[1]),
			Version:    strings.TrimSpace(requestLineStr[2]),
		},
		Headers: headers,
		Body:    []byte(payloadPart),
	}, nil
}

func (hr *HttpRequest) String() string {
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	buffer.WriteString(hr.RequestLine.HttpMethod + "\n")
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

func (hr *HttpRequest) PathParams() []string {
	params := strings.Split(hr.RequestLine.Target, "/")[1:] //ignore the first space
	length := len(params)
	if length > 0 && params[length-1] == "" {
		params = params[:length-1]
	}
	return params
}
