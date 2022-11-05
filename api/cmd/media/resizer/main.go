package main

import (
	"log"
	"os"

	cmd "github.com/and-period/furumaru/api/internal/media/cmd/resizer"
)

func main() {
	if err := cmd.Exec(); err != nil {
		log.Printf("An error has occurred: %v", err)
		os.Exit(1)
	}
}
