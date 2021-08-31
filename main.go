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

var version = "v0.0.0"

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	log.Printf("start debug-server %s", version)
	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("stopped")
}
