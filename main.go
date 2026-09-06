package main

import (
	"fmt"
	"os"
	"log"
)

func main() {
	f, err := os.Open("messages")
	if err != nil {
		log.Fatal("error:", err)
	}

	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", string(data[:n]))
	}
}
