package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"strings"
)

func TestRequestFromReader(t *testing.T) {

	goodRequests := []struct {
		inRequest string
		expectedParsed RequestLine
	}{
		{
			inRequest: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n", 
			expectedParsed: RequestLine{Method: "GET", RequestTarget: "/", HttpVersion: "1.1"},
		},
		{
			inRequest: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n", 
			expectedParsed: RequestLine{Method: "GET", RequestTarget: "/coffee", HttpVersion: "1.1"},
		},
	}

	badRequests := []string {
		"/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		"get /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		// "HELLO /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		"coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		" GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		"GET / http/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		"GET / HTTP1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		"GET / HTTP/1.1\rHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	}


	for i := range goodRequests {
		request := &goodRequests[i].inRequest
		expectedParsed := &goodRequests[i].expectedParsed

		r, err := RequestFromReader(strings.NewReader(*request))
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, expectedParsed, &r.RequestLine)
		// assert.Equal(t, expectedParsed.RequestLine.Method, r.RequestLine.Method)
		// assert.Equal(t, expectedParsed.RequestLine.RequestTarget, r.RequestLine.RequestTarget)
		// assert.Equal(t, expectedParsed.RequestLine.HttpVersion, r.RequestLine.HttpVersion)
	}

	for i, test := range badRequests {
		_, err := RequestFromReader(strings.NewReader(test))
		require.Error(t, err, "Bad Test %v, Input:\n%v", i, test)
	}
}
