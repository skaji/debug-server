package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	expectHeader := os.Getenv("EXPECT_HEADER")
	if expectHeader == "" {
		fmt.Fprintln(os.Stderr, "Need EXPECT_HEADER")
		os.Exit(1)
	}

	server := &Server{
		Server: &http.Server{
			Handler: &Handler{ExpectHeader: expectHeader},
			Addr:    ":8080",
		},
		WaitBeforeStop: 5 * time.Second,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		log.Printf("catch signal %s", <-quit)
		cancel()
	}()
	log.Println("start")
	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("stopped")
}
