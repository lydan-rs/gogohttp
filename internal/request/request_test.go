package request

import (
	"fmt"
	"http-protocol/internal/headers"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Include tests for nil values.

func testGoodRequestLines(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		requestLine RequestLine
		length      int
	}{
		{
			name:        "Simple GET",
			input:       "GET / HTTP/1.1\r\n",
			requestLine: RequestLine{Method: "GET", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("GET / HTTP/1.1\r\n"),
		},
		{
			name:        "GET /coffee",
			input:       "GET /coffee HTTP/1.1\r\n",
			requestLine: RequestLine{Method: "GET", RequestTarget: "/coffee", HttpVersion: "1.1"},
			length:      len("GET /coffee HTTP/1.1\r\n"),
		},
		{
			name:        "Simpe POST",
			input:       "POST / HTTP/1.1\r\n",
			requestLine: RequestLine{Method: "POST", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("POST / HTTP/1.1\r\n"),
		},
		{
			name:        "Simple PUT",
			input:       "PUT / HTTP/1.1\r\n",
			requestLine: RequestLine{Method: "PUT", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("PUT / HTTP/1.1\r\n"),
		},
		{
			name:        "Simple DELETE",
			input:       "DELETE / HTTP/1.1\r\n",
			requestLine: RequestLine{Method: "DELETE", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("DELETE / HTTP/1.1\r\n"),
		},
		{
			name:        "GET with headers",
			input:       "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			requestLine: RequestLine{Method: "GET", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("GET / HTTP/1.1\r\n"),
		},
		{
			name:        "GET /coffee with headers",
			input:       "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			requestLine: RequestLine{Method: "GET", RequestTarget: "/coffee", HttpVersion: "1.1"},
			length:      len("GET /coffee HTTP/1.1\r\n"),
		},
		{
			name:        "POST with headers",
			input:       "POST / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			requestLine: RequestLine{Method: "POST", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("POST / HTTP/1.1\r\n"),
		},
		{
			name:        "PUT with headers",
			input:       "PUT / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			requestLine: RequestLine{Method: "PUT", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("PUT / HTTP/1.1\r\n"),
		},
		{
			name:        "DELETE with headers",
			input:       "DELETE / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			requestLine: RequestLine{Method: "DELETE", RequestTarget: "/", HttpVersion: "1.1"},
			length:      len("DELETE / HTTP/1.1\r\n"),
		},
	}

	for i := range tests {
		request := tests[i].input
		expectedRequestLine := tests[i].requestLine
		expectedNumParsed := tests[i].length

		requestLine, numParsed, err := parseRequestLine([]byte(request))

		failureMsg := fmt.Sprintf("Subset: Good request lines.\nName: %v\nRequest: %v\n", tests[i].name, request)

		require.NoError(t, err, failureMsg)
		require.NotNil(t, requestLine, failureMsg)
		assert.Equal(t, expectedRequestLine, *requestLine, failureMsg)
		assert.Equal(t, expectedNumParsed, numParsed, failureMsg)
	}
}

func testIncompleteRequestLines(t *testing.T) {
	tests := []string{
		"",
		"GET",
		"GET / HT",
		"GET / HTTP/1.1",
	}

	for i := range tests {
		request := tests[i]
		requestLine, numParsed, err := parseRequestLine([]byte(request))

		failureMsg := fmt.Sprintf("Subset: Incomplete request lines.\nRequest: %v\n", request)

		require.NoError(t, err, failureMsg)
		assert.Nil(t, requestLine, failureMsg)
		assert.Equal(t, numParsed, 0, failureMsg)
	}
}

func testBadRequestLines(t *testing.T) {
	tests := []struct {
		name    string
		request string
		length  int
	}{
		{
			name:    "Missing method",
			request: "/coffee HTTP/1.1\r\n",
			length:  len("/coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Lowercase method",
			request: "get /coffee HTTP/1.1\r\n",
			length:  len("get /coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Bad method",
			request: "HELLO /coffee HTTP/1.1\r\n",
			length:  len("HELLO /coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Missing method bad resource path",
			request: "coffee HTTP/1.1\r\n",
			length:  len("coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Leading whitespace",
			request: " GET /coffee HTTP/1.1\r\n",
			length:  len(" GET /coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Illformed http version",
			request: "GET / http/1.1\r\n",
			length:  len("GET / http/1.1\r\n"),
		},
		{
			name:    "Illformed http version (2)",
			request: "GET / HTTP1.1\r\n",
			length:  len("GET / HTTP1.1\r\n"),
		},
		{
			name:    "Missing method with headers",
			request: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("/coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Lowercase method with headers",
			request: "get /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("get /coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Bad method with headers",
			request: "HELLO /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("HELLO /coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Missing method bad resource path with headers",
			request: "coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Leading whitespace with headers",
			request: " GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len(" GET /coffee HTTP/1.1\r\n"),
		},
		{
			name:    "Illformed http version with headers",
			request: "GET / http/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("GET / http/1.1\r\n"),
		},
		{
			name:    "Illformed http version with headers (2)",
			request: "GET / HTTP1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("GET / HTTP1.1\r\n"),
		},
		{
			name:    "Illformed CRLF",
			request: "GET / HTTP/1.1\rHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			length:  len("GET / HTTP/1.1\rHost: localhost:42069\r\n"),
		},
	}

	for i := range tests {
		request := tests[i].request
		expectedNumParsed := tests[i].length
		requestLine, numParsed, err := parseRequestLine([]byte(request))

		failureMsg := fmt.Sprintf("Subset: Bad request lines.\nName: %v\nRequest: %v\n", tests[i].name, request)

		require.Error(t, err, failureMsg)
		assert.Nil(t, requestLine, failureMsg)
		assert.Equal(t, expectedNumParsed, numParsed, failureMsg)
	}

}

func TestParseRequestLine(t *testing.T) {
	testGoodRequestLines(t)
	testIncompleteRequestLines(t)
	testBadRequestLines(t)
}

// ---- SECTION: parseBody ----

func TestParseBody(t *testing.T) {
	failureMsg := func(desc string) string { return fmt.Sprintf("Description: %v\n", desc) }
	

	desc := "Content length matches content-length. Expect body containing \"Hello There!\\n\"."
	content := "Hello There!\n"
	header := headers.Headers{
		"content-length": strconv.Itoa(len(content)),
	}
	body, numBytesParsed, finished, err := parseBody([]byte(content), header)
	require.NoError(t, err, failureMsg(desc))
	require.NotNil(t, body, failureMsg(desc))
	assert.Equal(t, content, string(body), failureMsg(desc))
	assert.Equal(t, len(content), numBytesParsed, failureMsg(desc))
	assert.True(t, finished, failureMsg(desc))


	desc = "Content length exceeds content-length. Expect body containing \"Hello There!\\n\"."
	content = "Hello There!\n General Kenobi!\n"
	target := "Hello There!\n"
	header = headers.Headers{
		"content-length": strconv.Itoa(len(target)),
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.NoError(t, err, failureMsg(desc))
	require.NotNil(t, body, failureMsg(desc))
	assert.Equal(t, target, string(body), failureMsg(desc))
	assert.Equal(t, len(target), numBytesParsed, failureMsg(desc))
	assert.True(t, finished, failureMsg(desc))


	desc = "Content length does not meet content-length. Expect nil body and no errors."
	content = "Hello The"
	target = "Hello There!\n"
	header = headers.Headers{
		"content-length": strconv.Itoa(len(target)),
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.NoError(t, err, failureMsg(desc))
	require.Nil(t, body, failureMsg(desc))
	assert.Equal(t, 0, numBytesParsed, failureMsg(desc))
	assert.False(t, finished, failureMsg(desc))


	desc = "Content empty. Expect nil body and no errors."
	content = ""
	target = "Hello There!\n"
	header = headers.Headers{
		"content-length": strconv.Itoa(len(target)),
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.NoError(t, err, failureMsg(desc))
	require.Nil(t, body, failureMsg(desc))
	assert.Equal(t, 0, numBytesParsed, failureMsg(desc))
	assert.False(t, finished, failureMsg(desc))


	desc = "No content-length header. Expect nil body, finished and no errors."
	content = ""
	header = headers.Headers{
		"some-name": "some-value",
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.NoError(t, err, failureMsg(desc))
	require.Nil(t, body, failureMsg(desc))
	assert.Equal(t, 0, numBytesParsed, failureMsg(desc))
	assert.True(t, finished, failureMsg(desc))


	desc = "Body content and no content-length header. Expect nil body, finished and no errors."
	content = "Some body content I guess."
	header = headers.Headers{
		"some-name": "some-value",
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.NoError(t, err, failureMsg(desc))
	require.Nil(t, body, failureMsg(desc))
	assert.Equal(t, 0, numBytesParsed, failureMsg(desc))
	assert.True(t, finished, failureMsg(desc))


	desc = "content-length value is a string. Expect error."
	content = "Some body content I guess."
	header = headers.Headers{
		"content-length": "some-value",
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.Error(t, err, failureMsg(desc))
	require.Nil(t, body, failureMsg(desc))
	require.Equal(t, 0, numBytesParsed, failureMsg(desc))
	require.False(t, finished, failureMsg(desc))


	desc = "content-length value is negative. Expect Errors"
	content = "Some body content I guess."
	header = headers.Headers{
		"content-length": "-3",
	}
	body, numBytesParsed, finished, err = parseBody([]byte(content), header)
	require.Error(t, err, failureMsg(desc))
	require.Nil(t, body, failureMsg(desc))
	require.Equal(t, 0, numBytesParsed, failureMsg(desc))
	require.False(t, finished, failureMsg(desc))
}

// ---- SECTION: Request From Reader ----

func TestRequestFromReader(t *testing.T) {

	// Test: Standard Body
	reader := strings.NewReader(
		"POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
	)
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))

	// Test: Body shorter than reported content length
	reader = strings.NewReader(
		"POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
	)
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}
