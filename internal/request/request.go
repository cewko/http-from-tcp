package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported HTTP version")
var SEPARATOR = "\r\n"

func ParseRequestLine(str string) (*RequestLine, string, error) {
	index := strings.Index(str, SEPARATOR)
	if index == -1 {
		return nil, str, nil
	}

	startLine := str[:index]
	remainingMessage := str[index+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, remainingMessage, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, remainingMessage, ERROR_MALFORMED_REQUEST_LINE
	}

	requestLine := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return requestLine, remainingMessage, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll: %w", err))
	}

	data_str := string(data)
	requestLine, _, err := ParseRequestLine(data_str)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, err
}
