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
		if lb > -1 {
			fmt.Printf("read: %s\n", lineBuilder.String())
			lineBuilder.Reset()
			// To start building the next line, we need to start at one position
			// past the new line character.
			// TODO: The way we are doing things assumes that in any chunk only one newline will be present.
			// This works for now, but with other input could easily break. Probably need a loop to
			// continually recontstruct new lines until there is no newline detected.
			constructLineFromBuf(buf[lb+1:nBytes], nBytes-(lb+1), &lineBuilder)
		}

	}

	if lineBuilder.Len() > 0 {
		fmt.Printf("read: %s\n", lineBuilder.String())
	}
}
