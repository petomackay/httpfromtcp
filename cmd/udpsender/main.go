package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal("Coult't resolve the address!")
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Couldn't dial the connection!")
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("unable to read input")
		}
		n, err := conn.Write([]byte(input))
		if err != nil {
			log.Fatalf("trouble writing to UDP: %e", err)
		}
		fmt.Printf("sent %d bytes\n", n)

	}

}
