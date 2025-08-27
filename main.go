package main

import (
    "os"
    "log"
    "io"
    "fmt"
)

func main() {
    file, err := os.Open("messages.txt")
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    buf := make([]byte, 8)

    for {
        nBytes, err := file.Read(buf)

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("read: %s\n", buf[:nBytes])
    }
}
