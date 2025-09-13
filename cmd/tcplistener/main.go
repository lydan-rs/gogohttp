package main

import (
	"net"
	"log"
	"fmt"
	"http-protocol/internal/request"
)



func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	fmt.Printf("TCP Listener started at %s.\n", listener.Addr())

	if err != nil { log.Fatal(err) }
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil { log.Fatal(err) }

		if addr := conn.RemoteAddr(); addr != nil{
			fmt.Printf("Connection Accepted from %s.\n", addr.String())
		} else {
			fmt.Printf("Connection Accepted from unkown source.\n")
		}

		request, err := request.RequestFromReader(conn)
		if err != nil { log.Fatal(err) }
		fmt.Println("Finished Reading Request.")

		response := fmt.Sprintf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion,
		)

		fmt.Printf("Response >>\n%v\n", response)
		conn.Write([]byte(response))
		fmt.Println("Response sent.")
	}
}
