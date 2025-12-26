package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine  RequestLine
	requestState requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const bufferSize = 8

var httpMethods = []string{"GET", "POST"}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		// we don't have enough data yet
		return 0, nil, nil
	}

	rlString := string(data[:idx])
	segments := strings.Split(rlString, " ")

	if len(segments) != 3 {
		return 0, nil, errors.New("Invalid number of segments in the request line")
	}

	httpMethodSegment := segments[0]
	httpTargetSegment := segments[1]
	httpProtocolSegment := segments[2]
	if !slices.Contains(httpMethods, httpMethodSegment) {
		return 0, nil, fmt.Errorf("Invalid HTTP method: %s", httpMethodSegment)
	}
	httpProtocolParts := strings.Split(httpProtocolSegment, "/")
	if httpProtocolParts[0] != "HTTP" {
		return 0, nil, fmt.Errorf("Invalid protocol name. Must be HTTP, was: %s", httpProtocolParts[0])
	}

	if httpProtocolParts[1] != "1.1" {
		return 0, nil, fmt.Errorf("Unsupported protocol version: %s", httpProtocolParts[1])
	}

	requestLine := &RequestLine{
		HttpVersion:   httpProtocolParts[1],
		RequestTarget: httpTargetSegment,
		Method:        httpMethodSegment,
	}

	// return idx + 2 to account for the CRLF bytes
	return idx + 2, requestLine, nil
}

func (r *Request) parse(data []byte) (int, error) {
	// TODO: switch instead
	if r.requestState == 0 {
		n, rl, err := parseRequestLine(data)
		if err != nil || n == 0 {
			// either there was an error, or no error == need more data
			return 0, err
		}
		r.RequestLine = *rl
		r.requestState = 1
		return n, nil
	}
	if r.requestState == 1 {
		return 0, fmt.Errorf("Trying to read data in a done state")
	}
	return 0, fmt.Errorf("Unkown request state: %d", r.requestState)
}

func RequestFromReader(request io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		requestState: requestStateInitialized,
	}

	for req.requestState != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err := request.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.requestState = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += n
		nParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[nParsed:])
		readToIndex -= nParsed
	}

	return req, nil
}
