package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		data := make([]byte, 8)
		line := ""
		for {
			count, err := f.Read(data)
			if err != nil {
				if err == io.EOF {
					if len(line) != 0 {
						ch <- line
						line = ""
					}
					return
				}
				fmt.Printf("fuckerror: %e", err)
			}
			parts := strings.Split(string(data[:count]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				line += parts[i]
				ch <- line
				line = ""
			}
			line += parts[len(parts)-1]
		}
	}()
	return ch
}

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error when accepting connection: %e", err)
		}
		fmt.Println("Connection accepted")
		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection closed")
	}
}
