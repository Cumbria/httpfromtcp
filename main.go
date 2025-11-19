package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
    ln, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("could not start listener: %v\n", err)
    }
    defer ln.Close()

    fmt.Println("Listening for TCP traffic on", port)
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("accept error: %v\n", err)
            continue
        }

        addr := conn.RemoteAddr().String()
        fmt.Printf("Accepted connection from %s\n", addr)

        linesChan := getLinesChannel(conn)
        for line := range linesChan {
            // print each line with a terminating newline, no extra formatting
            fmt.Println(line)
        }

        fmt.Printf("Connection from %s closed\n", addr)
    }
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
