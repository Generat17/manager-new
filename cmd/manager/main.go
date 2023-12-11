package main

import (
	"context"
	"log"
	"manager/internal/handler"
	"manager/internal/server"
	"os/signal"
	"sync"
	"syscall"

	"manager/internal/repository"
	"manager/internal/service"
	"manager/pkg/config"
)

const (
	configPath = "config.yaml"
)

func main() {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("failed to init config: %s", err.Error())
	}

	repo := repository.New()
	s := service.New(repo, cfg.FilePath, cfg.RecordTypes)

	// init storage
	err = s.WriteStorageFromFile()
	if err != nil {
		log.Printf("Failed to read storage file: %s", err.Error())
	}

	h := handler.New(s)
	serv := server.New(h, cfg.ServerPort)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	serv.Run(ctx, wg)

	<-ctx.Done()
	wg.Wait()

	err = s.UpdateFile()
	if err != nil {
		log.Fatalf("failed to update file: %s", err.Error())
	}

	log.Print("successful completion")
}
