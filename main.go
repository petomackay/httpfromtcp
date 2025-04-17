package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Couldn't start listener on port 42069: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Couln't accept connection: %v\n", err)
		}
		fmt.Println("A connection has been accepted")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("The connection has been closed.")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		buf := make([]byte, 8)
		current_line := ""
		for {
			read, err := f.Read(buf)
			if err != nil {
				if current_line != "" {
					// dump any leftovers
					lines <- current_line
				}
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("error reading the input file: %v", err)
			}
			chunk := string(buf[:read])
			parts := strings.Split(chunk, "\n")
			// min lenth should always be 1, so this is safe
			for _, part := range parts[:len(parts)-1] {
				lines <- current_line + part
				current_line = ""
			}
			current_line += parts[len(parts)-1]
		}
	}()
	return lines
}
