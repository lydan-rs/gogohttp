package main

import (
	"os"
	"log"
	"io"
	"fmt"
	"slices"
	"strings"
)

/*
Appends the contents of the buffer to the line builder.
Will only append up to the next newline character (if it exists).
Returns the index of the next newline character, or -1 if none is present.
*/
func constructLineFromBuf(buf []byte, nBytes int, builder *strings.Builder) int {
	linebreak := slices.Index(buf, '\n')
	if linebreak > -1 {
		builder.Write(buf[:linebreak])
	} else {
		builder.Write(buf[:nBytes])
	}

	return linebreak
}

func main() {
	file, err := os.Open("messages.txt")
	// file, err := os.Open("messages_multi_newline.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	buf := make([]byte, 8)
	var lineBuilder strings.Builder

	for {
		nBytes, err := file.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		lb := constructLineFromBuf(buf, nBytes, &lineBuilder)

		// If a newline is detected, then we need to print, flush, and restart the linebuilder.
		// For loop handles to possibility of multiple newlines in a single chunk.
		// When subslicing, the Index function will return a value relative to the start of the
		// subslice, not the base slice. So we need this offset to track where in the original slice
		// we are working from. Probably a better way to do this but for now its fine.
		offset := 0
		for lb > -1 {
			fmt.Printf("read: %s\n", lineBuilder.String())
			lineBuilder.Reset()
			offset += lb+1
			lb = constructLineFromBuf(buf[offset:nBytes], nBytes-offset, &lineBuilder)
		}

	}

	if lineBuilder.Len() > 0 {
		fmt.Printf("read: %s\n", lineBuilder.String())
	}
}
