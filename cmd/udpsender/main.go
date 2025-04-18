package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = conn.Write(line)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
