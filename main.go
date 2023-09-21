package main

import (
	"T/app"
	"T/services"
	"T/tools/config"
	"T/transport"
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	// read config from env
	cfg, err := config.Read()
	if err != nil {
		return errors.Wrap(err, "Config")
	}

	args := os.Args[1:]
	if len(args) < 1 {
		return errors.Errorf("Enter command line arg N. Example (N=5): ./program 5")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Errorf("Cannot parse command line arg N")
	}

	// init logs
	logger, err := app.InitLogs()
	if err != nil {
		return errors.Wrap(err, "InitLogs")
	}

	// domain Impl
	// ToDo: N param
	service := services.NewQueueService(logger, n)
	service.RunBackground()

	t := transport.NewQueueTransportService(logger, service)

	// Router
	router := http.NewServeMux()
	router.HandleFunc("/reg", t.AddTask)
	router.HandleFunc("/get", t.GetTasks)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Program ended!")

	return nil
}
