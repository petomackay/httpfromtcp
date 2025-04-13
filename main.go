package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Couldn't open the file: %v", err)
	}
	defer f.Close()
	buf := make([]byte, 8)

	var current_line string
	for {
		read, err := f.Read(buf)
		if err != nil {
			if current_line != "" {
				// dump any lelftovers
				fmt.Printf("read: %s\n", current_line)
			}
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("error reading the input file: %v", err)
		}
		if read == 0 {
			break
		}
		chunk := string(buf[:read])
		parts := strings.Split(chunk, "\n")
                if len(parts) > 1 {
			// there was a newline in the chunk
			for _, part := range parts[:len(parts) - 1] {
			    line := current_line + part
		            fmt.Printf("read: %s\n", line)
			    current_line = ""
			}
		}
		current_line += parts[len(parts) - 1]
	}
}
