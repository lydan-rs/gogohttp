package server

import (
	"fmt"
	"http-protocol/internal/request"
	"http-protocol/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

const crlf = "\r\n"

type Handler func(w response.Writer, req *request.Request)

type Server struct {
	tcpListener net.Listener
	connections []net.Conn
	handler Handler
	open atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil { return nil, err }

	connections := make([]net.Conn, 0, 32)
	
	server := Server{
		tcpListener: tcpListener,
		connections: connections,
		handler: handler,
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

	for {
		request, err := request.RequestFromReader(conn)
		writer := response.MakeWriter(conn)
		if request != nil {
			s.handler(writer, request)
			log.Printf("Response to %v sent.\nRequest: %v\n", conn.RemoteAddr().String(), request)
		}
		

		if err != nil {
			if err == io.EOF {
				log.Printf("%v has terminated connection.\n\n", conn.RemoteAddr().String())
			} else {
				log.Printf("---- ERROR ----\n>> Client: %v\n>> Error: %v\n\n", conn.RemoteAddr().String(), err.Error())
			}
			break
		}
	}
}
