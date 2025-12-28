package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	h := make(map[string]string)
	return h
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte("\r\n"))
	if crlfIdx == -1 {
		// incomplete header value
		return 0, false, nil
	}
	if crlfIdx == 0 {
		return 2, true, nil
	}

	hlString := strings.TrimSpace(string(data[:crlfIdx]))
	fmt.Printf("the header line is: '%s'\n", hlString)

	separatorIndex := strings.Index(hlString, ":")
	if separatorIndex == -1 {
		err := errors.New("Invalid header line format: no ':'")
		return 0, false, err
	}

	header := hlString[:separatorIndex]
	if len(header) != len(strings.TrimSpace(header)) {
		err := fmt.Errorf("Header field value contains whitespace: '%s'", header)
		return 0, false, err
	}

	value := strings.TrimSpace(hlString[separatorIndex+1:])
	fmt.Printf("the header: '%s'\nhas value: '%s'\n", header, value)

	h[header] = value

	return crlfIdx + 2, false, nil
}
