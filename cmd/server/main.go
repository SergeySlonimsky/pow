package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SergeySlonimsky/pow/internal/server"
	"github.com/SergeySlonimsky/pow/internal/server/cache"
	"github.com/SergeySlonimsky/pow/internal/server/handler"
	"github.com/SergeySlonimsky/pow/internal/server/pow"
	"github.com/SergeySlonimsky/pow/internal/server/storage"
)

func main() {
	ctx := context.Background()

	quoteStorage := storage.New()

	redis, err := cache.NewRedisCache(ctx, os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatalf("can't connect ro redis: %s", err.Error())
	}

	proofOfWork := pow.New(redis)

	h := handler.New(quoteStorage, proofOfWork)
	app := server.New(h)

	if err := app.Run(ctx, fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))); err != nil {
		log.Fatal(err)
	}
}
