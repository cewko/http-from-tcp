package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLine := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			for i := 0; i < n; i++ {
				if data[i] == '\n' {
					ch <- currentLine
					currentLine = ""
				} else {
					currentLine += string(data[i])
				}
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				break
			}
		}

		if currentLine != "" {
			ch <- currentLine
		}
	}()

	return ch
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}
