package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func main() {
	start := time.Now()
	fmt.Println("Start..")
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Printf("Done: %s\n", time.Since(start))
}

func run() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	ctx := context.Background()

	service, err := getService(ctx)
	if err != nil {
		return err
	}
	logger.Debug("Succeeded to get service")

	logger.Debug("Start to get channels")
	res, err := service.Channels.List([]string{"id", "snippet"}).Mine(true).Context(ctx).Do()
	if err != nil {
		logger.Error("Failed to get channels", zap.Error(err))
		return err
	}
	logger.Info("Succeeded to get channels", zap.Any("res", res))
	return nil
}
