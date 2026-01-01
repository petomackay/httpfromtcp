package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/petomackay/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine  RequestLine
	Headers      headers.Headers
	Body         []byte
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
	requestStateParsingHeaders
	requestStateParsingBody
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
	switch r.requestState {
	case requestStateInitialized:
		n, rl, err := parseRequestLine(data)
		if err != nil || n == 0 {
			// either there was an error, or no error == need more data
			return 0, err
		}
		r.RequestLine = *rl
		r.requestState = requestStateParsingHeaders
		return n, nil

	case requestStateParsingHeaders:
		h := r.Headers
		n, done, err := h.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			nextState := requestStateParsingBody
			if r.Headers.Get("content-length") == "" {
				nextState = requestStateDone
			}
			r.requestState = nextState
			//fmt.Println("Header parsing done")
		}
		return n, nil
	case requestStateParsingBody:
		contentLengthValue := r.Headers.Get("content-length")
		if contentLengthValue == "" {
			r.requestState = requestStateDone
			return 0, fmt.Errorf("Invalid state: parsing body with no content-length header!")
		}
		contentLength, err := strconv.Atoi(contentLengthValue)
		if err != nil {
			return 0, err
		}
		if r.Body == nil {
			r.Body = make([]byte, 0, contentLength)
		}
		//fmt.Println("parsing body with:")
		//fmt.Printf("BEGIN|%v|END", data)
		//fmt.Println("\nend of body")
		fmt.Printf("The length of body is: %d\n", len(data))
		fmt.Printf("The capacity of body is: %d\n", cap(data))
		r.Body = append(r.Body, data[:len(data)]...)
		fmt.Printf("the current body length: %d\ncontent length: %d\n", len(r.Body), contentLength)
		if len(r.Body) == contentLength {
			r.requestState = requestStateDone
		}
		if len(r.Body) > contentLength {
			return 0, fmt.Errorf("Body is longer than content-length")
		}
		return len(data), nil
	case requestStateDone:
		return 0, fmt.Errorf("Trying to read data in a done state")
	default:
		return 0, fmt.Errorf("Unkown request state: %d", r.requestState)
	}
}

func RequestFromReader(request io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		requestState: requestStateInitialized,
		Headers:      headers.NewHeaders(),
	}

	doneReading := false

	for req.requestState != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		var n int
		var err error
		if !doneReading {
			n, err = request.Read(buf[readToIndex:])
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				// we still might have data in the buffer that needs parsing
				doneReading = true
				if readToIndex == 0 {
					if req.requestState != requestStateDone {
						return nil, fmt.Errorf("Incomplete request")
					}
					break
				}
			} else {
				return nil, err
			}
		}
		readToIndex += n
		nParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if doneReading && nParsed == 0 {
			return nil, fmt.Errorf("Incomplete request")
		}
		copy(buf, buf[nParsed:])
		readToIndex -= nParsed
	}
	//fmt.Println("REQUEST PARSING DONE")

	return req, nil
}
