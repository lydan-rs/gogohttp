package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	

	// Test: Valid single header
	headers := MakeHeadersMap()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Valid single header without the final crlf
	headers = MakeHeadersMap()
	data = []byte("Host: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)


	// Test: Valid multiple values for the same name
	headers = MakeHeadersMap()
	data = []byte("Pets: Monte\r\nPets: Poppy\r\nPets: Bonnie\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "Monte, Poppy, Bonnie", headers["pets"])
	assert.Equal(t, 40, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = MakeHeadersMap()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 38, n)
	assert.False(t, done)

	// Test: Invalid field name character
	headers = MakeHeadersMap()
	data = []byte("<host>: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 25, n)
	assert.False(t, done)
}
