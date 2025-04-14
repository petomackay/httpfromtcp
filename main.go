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

        lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <- chan string {
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
			for _, part := range parts[:len(parts) - 1] {
				lines <- current_line + part
				current_line = ""
			}
			current_line += parts[len(parts) - 1]
		}
	}()
	return lines
}
