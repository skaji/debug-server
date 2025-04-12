package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{
		Handler: &Handler{},
		Addr:    ":8080",
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	done := make(chan error)
	go func() {
		fmt.Printf("listen http://localhost%s\n", server.Addr)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		done <- err
	}()

	err := func() error {
		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			fmt.Println("catch signal, shutdown server...")
			ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel2()
			if err := server.Shutdown(ctx2); err != nil {
				fmt.Println("failed to shutdown server, force to shutdown server")
				_ = server.Close()
			}
			return <-done
		}
	}()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("finished")
}

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		_, _ = w.Write([]byte("OK\n"))
		return
	}

	script, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(script) == 0 {
		http.Error(w, "need body", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bash", "-c", string(script)).CombinedOutput()
	if len(out) > 0 {
		_, _ = w.Write(out)
	}
	if err != nil {
		_, _ = fmt.Fprintln(w, err.Error())
	}
}
