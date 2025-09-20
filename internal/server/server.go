package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

const crlf = "\r\n"

type Server struct {
	tcpListener net.Listener
	connections []net.Conn
	open atomic.Bool
}

func Serve(port int) (*Server, error) {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil { return nil, err }

	connections := make([]net.Conn, 0, 32)
	
	server := Server{
		tcpListener: tcpListener,
		connections: connections,
	}
	server.open.Store(true)

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.open.Store(false)
	if s.tcpListener != nil {
			return s.tcpListener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			if !s.open.Load() {
				return
			}
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	n, err := conn.Write([]byte(
		"HTTP/1.1 200 OK" + crlf +
		"Content-Type: text/plain" + crlf +
		"Content-Length: 13" + crlf +
		crlf +
		"Hello World!\n"))

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%v bytes written to %v from server %v.\n", n, conn.RemoteAddr().String(), s.tcpListener.Addr().String())
	}
}
