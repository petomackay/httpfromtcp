package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Couldn't open the file:")
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	buf := make([]byte, 8)
	for {
	    read, err := f.Read(buf)
	    if err != nil {
		if errors.Is(err, io.EOF) {
		    os.Exit(0)
		}
		fmt.Println("error reading file:")
		fmt.Println(err)
		os.Exit(1)
	    }
	    if read == 0 {
		break
	    }
		fmt.Printf("read: %s\n", string(buf[:read]))
	}
}
