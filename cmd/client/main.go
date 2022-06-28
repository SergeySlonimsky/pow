package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SergeySlonimsky/pow/internal/client"
)

func main() {
	serverURL := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	if err := client.Run(context.Background(), serverURL); err != nil {
		log.Fatalf("run client: %s", err.Error())
	}
}
