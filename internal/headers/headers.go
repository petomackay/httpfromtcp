package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const specialChars string = "!#$%&'*+-.^_`|~"

type Headers map[string]string

func NewHeaders() Headers {
	h := make(map[string]string)
	return h
}

func (h Headers) Set(header, value string) {
	headerLC := strings.ToLower(header)
	val, ok := h[headerLC]
	if ok {
		value = fmt.Sprintf("%s, %s", val, value)
	}
	h[headerLC] = value
}

func (h Headers) Get(header string) string {
	headerLC := strings.ToLower(header)
	return h[headerLC]
}

func validateHeaderKey(header string) error {
	headerLength := len(header)

	if headerLength == 0 {
		return fmt.Errorf("Header key cannot be empty")
	}

	if headerLength != len(strings.TrimSpace(header)) {
		err := fmt.Errorf("Header field value contains whitespace: '%s'", header)
		return err
	}

	for pos, char := range header {
		if isValidTChar(char) {
			fmt.Printf("char: %c at position: %d is valid\n", char, pos)
		} else {
			fmt.Printf("char: %c at position: %d is invalid\n", char, pos)
			err := fmt.Errorf("Header key contains invalid char at position %d: %c", pos, char)
			return err
		}
	}

	return nil
}

func isValidTChar(c rune) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	if c >= '0' && c <= '9' {
		return true
	}
	if strings.ContainsRune(specialChars, c) {
		return true
	}
	return false
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
	err = validateHeaderKey(header)
	if err != nil {
		return 0, false, err
	}

	value := strings.TrimSpace(hlString[separatorIndex+1:])
	fmt.Printf("the header: '%s'\nhas value: '%s'\n", header, value)

	h.Set(header, value)

	return crlfIdx + 2, false, nil
}
