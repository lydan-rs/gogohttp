package main

import (
	"os"
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

	} ()

	return ch
}


func main() {
	file, err := os.Open("messages.txt")
	// file, err := os.Open("messages_multi_newline.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	lineCH := getLinesChannel(file)
	// I think `range` ends when the channel is closed. Otherwise I have no idea why this works.
	for line := range lineCH {
		fmt.Printf("read: %s\n", line)
	}

}
