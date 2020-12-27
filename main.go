package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	expectHeader := os.Getenv("EXPECT_HEADER")
	if expectHeader == "" {
		fmt.Fprintln(os.Stderr, "Need EXPECT_HEADER")
		os.Exit(1)
	}
	server := &http.Server{
		Handler: &Handler{ExpectHeader: expectHeader},
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
