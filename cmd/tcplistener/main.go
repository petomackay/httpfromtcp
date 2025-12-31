package main

import (
	"fmt"
	"log"
	"net"

	"github.com/petomackay/httpfromtcp/internal/request"
)

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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error when parsing request: %e", err)
		}
		requestLine := req.RequestLine
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", requestLine.Method)
		fmt.Printf("- Target: %s\n", requestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", requestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Connection closed")
	}
}
