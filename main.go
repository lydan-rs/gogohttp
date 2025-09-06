package main

import (
	"net"
	"log"
	"io"
	"fmt"
	"slices"
	"strings"
)

// <-chan is a receive only channel, meaning that data can only be pulled from it.
// The opposite would be a send only channel, chan<-
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() { 
		// Line Processing
		buf := make([]byte, 8)
		var lineBuilder strings.Builder

		for {
			nBytes, err := f.Read(buf)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}


			// If a newline is detected, then we need to print, flush, and restart the linebuilder.
			// For loop handles to possibility of multiple newlines in a single chunk.
			// When subslicing, the Index function will return a value relative to the start of the
			// subslice, not the base slice. So we need this offset to track where in the original slice
			// we are working from. Probably a better way to do this but for now its fine.
			offset := 0
			for {
				newlineIndex := slices.Index(buf[offset:nBytes], '\n')

				if newlineIndex > -1 {
					lineBuilder.Write(buf[offset:offset+newlineIndex])
					ch <- lineBuilder.String()
					lineBuilder.Reset()
					offset += newlineIndex+1
				} else {
					lineBuilder.Write(buf[offset:nBytes])
					break
				}
			}
		}

		if lineBuilder.Len() > 0 {
			ch <- lineBuilder.String()
		}

		close(ch)
		f.Close()
		fmt.Println("Finished reading from connection. Closing connection and channel.")

	} ()

	return ch
}


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

		lineCH := getLinesChannel(conn)
		// I think `range` ends when the channel is closed. Otherwise I have no idea why this works.
		for line := range lineCH {
			fmt.Printf("%s\n", line)
		}
	}
}
