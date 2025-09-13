package request

import (
	"bytes"
	"errors"
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

/*
HTTP-NAME: [prefix]HTTP
HTTP-VER: [HTTP-NAME]/[major digit].[minor digit]
REQUEST-LINE: [method] [resouce path] [HTTP-VER]
*/

func makeRequest() Request {
	r := Request{}
	r.state = initialised
	return r
}

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
		Method:        method,
		RequestTarget: resourcePath,
		HttpVersion:   "1.1",
	}, nil

}

func (r *Request) parse(line string) error {

	switch r.state {
	case initialised:
		requestLine, err := parseRequestLine(line)
		if err != nil {
			return err
		}
		r.RequestLine = *requestLine
		r.state = done
	}

	return nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 0, 1024)
	chunk := make([]byte, 8)
	lines := make([]string, 0, 16)
	nChunks := 0

	// TODO: Consider a version of this that allows you to parse a line as soon as it has been created.
	for {
		nChunks += 1
		fmt.Printf("Reading chunk %v from message.\n", nChunks)
		nRead, readErr := reader.Read(chunk)
		fmt.Printf("%v bytes read from connection.\n", nRead)
		buffer = append(buffer, chunk[:nRead]...)
		fmt.Println("Bytes appended to buffer.")
		fmt.Println()

		if nRead != 0 {
			/*
				If bytes were read, extract lines seperated by CRLF, until there are none.
				Copying the unprocessed bytes to the begining of the slice may not be the most efficient thing.
				But with how this loop works, you'll only be copying a max of len(chunk)-1 bytes each time.
				So for now its probably fine. Might want to check out in the future though.
			*/
			// TODO: Optimise by restricting the range to search for the CRLF to the bytes just read in, minus 1.
			/*
				Its possible that a CRLF gets split between two chunks. So to gurrantee that one is found, we would need
				to include the last byte appended before the current chunk in the search range to gurrantee an accurate find.
			*/
			for index := bytes.Index(buffer, []byte(crlf)); index >= 0; {
				lines = append(lines, string(buffer[:index]))
				copy(buffer, buffer[index+2:])
				buffer = buffer[:len(buffer)-(index+2)]
				index = bytes.Index(buffer, []byte(crlf))
			}
		}

		// HTTP messages don't necessarily end with a CRLF, so we need to make sure to grab all the remainging data on EOF.
		if readErr == io.EOF {
			lines = append(lines, string(buffer))
			fmt.Println("End Of File reached.")
			break
		}

		if readErr != nil {
			return nil, readErr
		}
	}

	fmt.Println("Connection Read Finished.")

	if len(lines) == 0 {
		return nil, errors.New("Message was empty.")
	}

	request := makeRequest()
	for lIdx := range lines {
		fmt.Println(lines[lIdx])
		err := request.parse(lines[lIdx])
		if err != nil {
			return nil, err
		}
	}
	// requestLine, err := parseRequestLine(lines[0])
	// if err != nil {
	// 	return nil, err
	// }

	return &request, nil

}
