package response

import (
	"fmt"
	"http-protocol/internal/headers"
	"io"
	"net"
	"strconv"
)

type StatusCode uint

const crlf = "\r\n"
const (
	SC_OK                    StatusCode = 200
	SC_BAD_REQUEST           StatusCode = 400
	SC_INTERNAL_SERVER_ERROR StatusCode = 500
)

func ReasonPhrase(code StatusCode) (string, error) {
	switch code {
	case SC_OK:
		return "OK", nil
	case SC_BAD_REQUEST:
		return "Bad Request", nil
	case SC_INTERNAL_SERVER_ERROR:
		return "Internal Server Error", nil
	default:
		return "", fmt.Errorf("%v is not a supported status code.", code)
	}
}

type writerStatus int

const (
	wS_INIT writerStatus = iota
	wS_STATUS_LINE
	wS_HEADERS
	wS_BODY
)

type Writer struct {
	connection net.Conn
	status     writerStatus
}

func MakeWriter(conn net.Conn) Writer {
	return Writer {
		connection: conn,
		status: wS_INIT,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.status != wS_INIT {
		return fmt.Errorf("Status line already writtern")
	}

	reasonPhrase, err := ReasonPhrase(statusCode)
	if err != nil {
		return err
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %v %v%v", statusCode, reasonPhrase, crlf)

	_, err = w.connection.Write([]byte(statusLine))
	w.status = wS_STATUS_LINE
	return err
}

func getDefaultHeaders(contentLen int) headers.Headers {
	h := headers.MakeHeadersMap()
	h.Add("Content-Length", strconv.Itoa(contentLen))
	h.Add("Connection", "close")
	h.Add("Content-Type", "text/plain")
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.status != wS_STATUS_LINE {
		w.WriteStatusLine(SC_OK)
	}

	for key := range headers {
		_, err := w.connection.Write([]byte(fmt.Sprintf("%v: %v%v", key, headers[key], crlf)))
		if err != nil {
			return err
		}
	}
	w.connection.Write([]byte(crlf))
	w.status = wS_HEADERS
	return nil
}

func (w *Writer) WriteString(body []byte) (int, error) {
	if w.status != wS_HEADERS {
		w.WriteHeaders(getDefaultHeaders(len(body)))
	}

	n, err := w.connection.Write(body)
	
	w.status = wS_BODY
	return n, err
}


func (w *Writer) WriteHTML(body []byte) (int, error) {
	if w.status != wS_HEADERS {
		h := getDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		w.WriteHeaders(h)
	}

	n, err := w.connection.Write(body)
	
	w.status = wS_BODY
	return n, err
}

func writeChunks(w io.Writer, p []byte, maxChunkSize uint) (int, error) {
	total := 0
	for i := 0; i < len(p); {
		chunkSize := min(int(maxChunkSize), len(p)-i)
		line := fmt.Sprintf("%x\r\n%s\r\n", chunkSize, p[i:i+chunkSize])
		n, err := w.Write([]byte(line))
		i += chunkSize
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func GetChunkHeaders(contentType string) headers.Headers {
	h := headers.MakeHeadersMap()
	h.Set("Content-Type", contentType)
	h.Set("Transfer-Encoding", "chunked")
	return h
}


func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.status != wS_HEADERS {
		w.WriteHeaders(GetChunkHeaders("text/plain"))
		w.status = wS_HEADERS
	}

	total := 0
	for total < len(p) {
		n, err := writeChunks(w.connection, p, 32)
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}


func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.status != wS_HEADERS {
		w.WriteHeaders(GetChunkHeaders("text/plain"))
		w.status = wS_HEADERS
	}

	n, err := w.connection.Write([]byte("0\r\n\r\n"))
	w.status = wS_BODY
	return n, err
}
