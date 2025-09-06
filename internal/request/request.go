package request

import (
	"io"
	"strings"
	"errors"
	"unicode"
	"fmt"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

/*
HTTP-NAME: [prefix]HTTP
HTTP-VER: [HTTP-NAME]/[major digit].[minor digit]
REQUEST-LINE: [method] [resouce path] [HTTP-VER]
*/

const CRLF string = "\r\n"



// TODO: Read up on RFCs to get a better idea of how to do this.
func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if nParts := len(parts); nParts != 3 {
		return nil, fmt.Errorf("Invalid Request Line. Contains %v parts, expected 3", nParts)
	}

	method := parts[0]
	resourcePath := parts[1]
	version := parts[2]
	
	// Method
	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return nil, fmt.Errorf("Invalid HTTP Method: %v", method)
		}
	}

	// Resource Path
	// TODO: Validate resouce path

	// HTTP Version
	if version != "HTTP/1.1" {
		return nil, fmt.Errorf("Invalid HTTP Version: %v", version) 
	}

	return &RequestLine{
		Method: method,
		RequestTarget: resourcePath,
		HttpVersion: "1.1",
	}, nil

}



// TODO: Reading message data and converting to lines - Do Better.
func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil { return nil, err }

	lines := strings.Split(string(buf), CRLF)
	if len(lines) == 0 {
		return nil, errors.New("Message was empty.")
	}
	
	requestLine, err := parseRequestLine(lines[0])
	if err != nil { return nil, err}

	return &Request{RequestLine: *requestLine}, nil

}
