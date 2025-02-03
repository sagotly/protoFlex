package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("/etc/wireguard/wg0.conf")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	fmt.Println(string(data))
}
