package main

import (
	"MusicLibrary/internal/api"
	"MusicLibrary/internal/config"
	"MusicLibrary/internal/database"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("creating config")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("initiating database")
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewServer(cfg, db)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Printf("starting server on %s:%d", cfg.Host.Address, cfg.Host.Port)
		return server.Start()
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Printf("closing database")
		return db.Close()
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Printf("shutting down server")
		return server.Stop()
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
