package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const crlf string = "\r\n"

type parsingState int

const (
	initialised = iota
	done
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parsingState
}


func makeRequest() Request {
	r := Request{}
	r.state = initialised
	return r
}


/*
HTTP-NAME: [prefix]HTTP
HTTP-VER: [HTTP-NAME]/[major digit].[minor digit]
REQUEST-LINE: [method] [resouce path] [HTTP-VER]
*/

// TODO: Read up on RFCs to get a better idea of how to do this.
func parseRequestLine(data []byte) (*RequestLine, int, error) {
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex < 0 {
		return nil, 0, nil
	}

	parts := strings.Split(string(data[:crlfIndex]), " ")
	bytesParsed := crlfIndex + 2

	if nParts := len(parts); nParts != 3 {
		return nil, bytesParsed, fmt.Errorf("Invalid Request Line. Contains %v parts, expected 3", nParts)
	}

	method := parts[0]
	resourcePath := parts[1]
	version := parts[2]

	// Method
	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return nil, bytesParsed, fmt.Errorf("Invalid HTTP Method: %v", method)
		}
	}

	// Resource Path
	// TODO: Validate resouce path

	// HTTP Version
	if version != "HTTP/1.1" {
		return nil, bytesParsed, fmt.Errorf("Invalid HTTP Version: %v", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: resourcePath,
		HttpVersion:   "1.1",
	}, bytesParsed, nil

}




func (r *Request) parse(data []byte) (int, error) {
	bytesConsumed := 0

	for r.state != done {
		switch r.state {
		case initialised:
			requestLine, bytesParsed, err := parseRequestLine(data)
			if err != nil || bytesParsed == 0 {
				return 0, err
			}
			fmt.Printf("RequestLine: %v\n", requestLine)
			r.RequestLine = *requestLine
			bytesConsumed += bytesParsed
			r.state = done
		}
	}

	return bytesConsumed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 0, 1024)
	chunk := make([]byte, 8)
	nChunks := 0

	request := makeRequest()

	for request.state != done {
		nChunks += 1
		fmt.Printf("Reading chunk %v from message.\n", nChunks)
		nRead, readErr := reader.Read(chunk)
		fmt.Printf("%v bytes read from connection.\n", nRead)
		buffer = append(buffer, chunk[:nRead]...)
		fmt.Println("Bytes appended to buffer.")
		fmt.Println()

		if nRead != 0 {
			bytesConsumed, err := request.parse(buffer)
			if err != nil {
				return nil, err
			}
			copy(buffer, buffer[bytesConsumed:])
			buffer = buffer[:len(buffer)-bytesConsumed]
		}

		if readErr != nil {
			if readErr == io.EOF {
				return nil, fmt.Errorf("Connection to ended prematurely.")
			} else {
				return nil, fmt.Errorf("Unexpected Error Expeirenced: %v.", readErr)
			}
		}
	}

	return &request, nil

}
